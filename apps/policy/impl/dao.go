package impl

import (
	"context"

	"github.com/zhou-lincong/keyauth/apps/policy"

	"github.com/infraboard/mcube/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Save Object
func (s *service) save(ctx context.Context, ins *policy.Policy) error {
	// s.col.InsertMany()
	if _, err := s.col.InsertOne(ctx, ins); err != nil {
		return exception.NewInternalServerError("inserted policy(%s) document error, %s",
			ins.Id, err)
	}
	return nil
}

func newQueryPolicyRequest(r *policy.QueryPolicyRequest) *queryRequest {
	return &queryRequest{
		r,
	}
}

// 把QueryReq --> MongoDB Options
type queryRequest struct {
	*policy.QueryPolicyRequest
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

	//这里没有做关键字查询

	if r.Namespace != "" {
		filter["spec.namespace"] = r.Namespace
	}
	if r.Username != "" {
		filter["spec.username"] = r.Username
	}
	if r.Role != "" {
		filter["spec.role"] = r.Role
	}

	return filter
}

// LIST, Query, 会很多条件(分页, 关键字, 条件过滤, 排序)
// 需要单独为其 做过滤参数构建
func (s *service) query(ctx context.Context, req *queryRequest) (*policy.PolicySet, error) {
	// SQL Where
	// FindFilter
	resp, err := s.col.Find(ctx, req.FindFilter(), req.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find policy error, error is %s", err)
	}

	set := policy.NewPolicySet()
	// 循环
	for resp.Next(ctx) {
		ins := policy.NewDefaultPolicy()
		if err := resp.Decode(ins); err != nil {
			return nil, exception.NewInternalServerError("decode policy error, error is %s", err)
		}

		set.Add(ins)
	}

	// count
	count, err := s.col.CountDocuments(ctx, req.FindFilter())
	if err != nil {
		return nil, exception.NewInternalServerError("get policy count error, error is %s", err)
	}
	set.Total = count

	return set, nil
}
