package impl

import (
	"context"

	"github.com/zhou-lincong/keyauth/apps/role"
)

func (s *service) CreateRole(ctx context.Context, req *role.CreateRoleRequest) (
	*role.Role, error) {
	ins := role.NewRole(req)

	if err := s.save(ctx, ins); err != nil {
		return nil, err
	}

	return ins, nil
}

func (s *service) QueryRole(ctx context.Context, req *role.QueryRoleRequest) (
	*role.RoleSet, error) {
	query := newQueryRequest(req)
	return s.query(ctx, query)
}
