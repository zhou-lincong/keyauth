package impl

import (
	"context"
	"fmt"

	"github.com/zhou-lincong/keyauth/apps/book"

	"github.com/infraboard/mcube/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *service) save(ctx context.Context, ins *book.Book) error {
	// s.col.InsertMany()插入多条数据
	if _, err := s.col.InsertOne(ctx, ins); err != nil {
		return exception.NewInternalServerError("inserted book(%s) document error, %s",
			ins.Data.Name, err)
	}
	return nil
}

// GET, 对应Describe
// filter 过滤器(Collection),类似于MYSQL Where条件
// 然后调用Decode方法来进行 反序列化  bytes ---> Object (通过BSON Tag)
func (s *service) get(ctx context.Context, id string) (*book.Book, error) {
	//条件只有一个，只需要传key的名字和字段
	//在表里面，想过滤哪个字段，filter就传哪个，例如user：admin，就没有sql注入的问题，直接把参数传过去
	filter := bson.M{"_id": id}
	ins := book.NewDefaultBook()
	if err := s.col.FindOne(ctx, filter).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("book %s not found", id)
		}

		return nil, exception.NewInternalServerError("find book %s error, %s", id, err)
	}

	return ins, nil
}

func newQueryBookRequest(r *book.QueryBookRequest) *queryBookRequest {
	return &queryBookRequest{
		r,
	}
}

// 方便条件的构建，把查询请求QueryReq -->转化成 MongoDB Options过滤条件
type queryBookRequest struct {
	*book.QueryBookRequest
}

// 做分页功能
func (r *queryBookRequest) FindOptions() *options.FindOptions {
	pageSize := int64(r.Page.PageSize)
	skip := int64(r.Page.PageSize) * int64(r.Page.PageNumber-1)

	opt := &options.FindOptions{
		// Value: -1 倒序--- 排序:用MySQL写就是 Order By create_at Desc
		Sort: bson.D{
			{Key: "create_at", Value: -1},
		},
		// 做分页需要用到的参数，分页: limit 0,10  skip:0, limit:10
		Limit: &pageSize,
		Skip:  &skip,
	}

	return opt
}

// 过滤条件
// 由于Mongodb支持嵌套, JSON, 如何过滤嵌套里面的条件, 使用.访问嵌套对象属性
func (r *queryBookRequest) FindFilter() bson.M {
	// mongo过滤的map
	filter := bson.M{}
	// 下面这种写法相当于添加了一个where的条件where key=value
	// filter["key"] = "value"
	// 过滤的语法，支持正则，但没有MySQL方便
	if r.Keywords != "" { //具体的过滤参数
		filter["$or"] = bson.A{
			bson.M{"data.name": bson.M{"$regex": r.Keywords, "$options": "im"}},
			bson.M{"data.author": bson.M{"$regex": r.Keywords, "$options": "im"}},
		}
	}
	return filter
}

// LIST, Query, 会很多条件(分页, 关键字, 条件过滤, 排序)
// 需要单独为其 做过滤参数构建
func (s *service) query(ctx context.Context, req *queryBookRequest) (*book.BookSet, error) {
	// 在MYSQL是通过 Where
	// 在mongodb是通过过滤器FindFilter
	resp, err := s.col.Find(ctx, req.FindFilter(), req.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find book error, error is %s", err)
	}

	set := book.NewBookSet()
	// 循环
	for resp.Next(ctx) {
		ins := book.NewDefaultBook()
		if err := resp.Decode(ins); err != nil {
			return nil, exception.NewInternalServerError("decode book error, error is %s", err)
		}

		set.Add(ins)
	}

	// count
	count, err := s.col.CountDocuments(ctx, req.FindFilter())
	if err != nil {
		return nil, exception.NewInternalServerError("get book count error, error is %s", err)
	}
	set.Total = count

	return set, nil
}

// UpdateByID, 通过主键来更新对象
func (s *service) update(ctx context.Context, ins *book.Book) error {
	// 等价MYSQL update obj(SET f=v,f=v) where id=?
	// s.col.UpdateOne(ctx, filter(), ins)主键不确定用这个
	if _, err := s.col.UpdateByID(ctx, ins.Id, ins); err != nil {
		return exception.NewInternalServerError("inserted book(%s) document error, %s",
			ins.Data.Name, err)
	}

	return nil
}

func (s *service) deleteBook(ctx context.Context, ins *book.Book) error {
	if ins == nil || ins.Id == "" {
		return fmt.Errorf("book is nil")
	}
	//等价 delete from book where id = ?
	result, err := s.col.DeleteOne(ctx, bson.M{"_id": ins.Id})
	if err != nil {
		return exception.NewInternalServerError("delete book(%s) error, %s", ins.Id, err)
	}

	if result.DeletedCount == 0 {
		return exception.NewNotFound("book %s not found", ins.Id)
	}

	return nil
}
