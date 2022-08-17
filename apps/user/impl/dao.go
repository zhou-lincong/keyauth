package impl

import (
	"context"
	"fmt"

	"github.com/zhou-lincong/keyauth/apps/user"

	"github.com/infraboard/mcube/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newQueryRequest(r *user.QueryUserRequest) *queryRequest {
	return &queryRequest{
		r,
	}
}

// 把QueryReq --> MongoDB Options
type queryRequest struct {
	*user.QueryUserRequest
}

// Find参数
func (r *queryRequest) FindOptions() *options.FindOptions {
	pageSize := int64(r.Page.PageSize)
	skip := int64(r.Page.PageSize) * int64(r.Page.PageNumber-1)

	opt := &options.FindOptions{
		// 排序: Order By create_at Desc
		Sort: bson.D{
			{Key: "create_at", Value: -1},
		},
		// 分页: limit 0,10  skip:0, limit:10
		Limit: &pageSize,
		Skip:  &skip,
	}

	return opt
}

// 过滤条件
// 由于Mongodb支持嵌套, JSON, 如何过滤嵌套里面的条件, 使用.访问嵌套对象属性
func (r *queryRequest) FindFilter() bson.M {
	filter := bson.M{}
	// where key=value
	// filter["key"] = "value"
	return filter
}

// LIST, Query, 会很多条件(分页, 关键字, 条件过滤, 排序)
// 需要单独为其 做过滤参数构建
func (s *impl) query(ctx context.Context, req *queryRequest) (*user.UserSet, error) {
	// SQL Where
	// FindFilter
	resp, err := s.col.Find(ctx, req.FindFilter(), req.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find book error, error is %s", err)
	}

	set := user.NewUserSet()
	// 循环
	for resp.Next(ctx) {
		ins := user.NewDefaultUser()
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

// GET, Describe
// filter 过滤器(Collection),类似于MYSQL Where条件
// 调用Decode方法来进行 反序列化  bytes ---> Object (通过BSON Tag)
func (i *impl) get(ctx context.Context, req *user.DescribeUserRequest) (*user.User, error) {
	filter := bson.M{}
	switch req.DescribeBy {
	case user.DescribeBy_USER_ID:
		filter["_id"] = req.UserId
	case user.DescribeBy_USER_NAME:
		filter["data.domain"] = req.Domain
		filter["data.name"] = req.UserName
	default:
		return nil, fmt.Errorf("unknow describe_by %s", req.DescribeBy)
	}

	ins := user.NewDefaultUser()
	if err := i.col.FindOne(ctx, filter).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("user %s not found", req)
		}

		return nil, exception.NewInternalServerError("find user %s error, %s", req, err)
	}

	return ins, nil
}
