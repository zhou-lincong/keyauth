package impl

import (
	"context"

	"github.com/zhou-lincong/keyauth/apps/policy"
	"github.com/zhou-lincong/keyauth/apps/role"

	"github.com/infraboard/mcube/exception"
)

//校验权限
func (s *service) ValidatePermission(ctx context.Context, req *policy.ValidatePermissionRequest) (
	*policy.Policy, error) {
	// 根据用户和命名空间找到该用户的授权策略
	// 由于使用分页, 只查询100条数据
	query := policy.NewQueryPolicyRequest()
	query.Namespace = req.Namespace
	query.Username = req.Username
	query.Page.PageSize = 100
	set, err := s.QueryPolicy(ctx, query)
	if err != nil {
		return nil, err
	}

	// 获取用户的角色, 从策略中抽取出来
	roleNames := set.Roles()

	// 通过Role模块查询所有的Role对象
	queryRoleReq := role.NewQueryRoleRequestWithName(roleNames)
	queryRoleReq.Page.PageSize = 100
	roles, err := s.role.QueryRole(ctx, queryRoleReq)
	if err != nil {
		return nil, err
	}

	// 根据Role判断用户是否具有权限
	hasPerm, role := roles.HasPermission(role.NewPermissionRequest(req.Service, req.Resource, req.Action))
	if !hasPerm {
		return nil, exception.NewPermissionDeny("not permission access service %s resource %s action %s",
			req.Service,
			req.Resource,
			req.Action,
		)
	}

	p := set.GetPolicyByRole(role.Spec.Name)
	return p, nil
}

func (s *service) QueryPolicy(ctx context.Context, req *policy.QueryPolicyRequest) (
	*policy.PolicySet, error) {
	query := newQueryPolicyRequest(req)
	return s.query(ctx, query)
}

func (s *service) CreatePolicy(ctx context.Context, req *policy.CreatePolicyRequest) (
	*policy.Policy, error) {
	ins, err := policy.NewPolicy(req)
	if err != nil {
		return nil, exception.NewBadRequest("validate create book error, %s", err)
	}

	if err := s.save(ctx, ins); err != nil {
		return nil, err
	}

	return ins, err
}
