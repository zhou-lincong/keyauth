package impl

import (
	"context"

	"github.com/zhou-lincong/keyauth/apps/endpoint"

	"github.com/infraboard/mcube/exception"
)

// Save Object
func (s *service) save(ctx context.Context, set *endpoint.EndpointSet) error {
	// s.col.InsertMany()	插入多个文档，需要传一个文档的列表
	if _, err := s.col.InsertMany(ctx, set.ToDocs()); err != nil {
		return exception.NewInternalServerError("inserted service %s endpoint document error, %s",
			set.Service, err)
	}
	return nil
}
