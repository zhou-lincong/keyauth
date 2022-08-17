package role

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	request "github.com/infraboard/mcube/http/request"
	"github.com/rs/xid"
)

const (
	AppName = "role"
)

var (
	validate = validator.New()
)

type Service interface {
	CreateRole(context.Context, *CreateRoleRequest) (*Role, error)
	RPCServer
}

func NewQueryRoleRequestWithName(names []string) *QueryRoleRequest {
	return &QueryRoleRequest{
		Page:      request.NewDefaultPageRequest(),
		RoleNames: names,
	}
}

//判断一组角色是不是有权限
func (s *RoleSet) HasPermission(req *PermissionRequest) (bool, *Role) {
	for i := range s.Items {
		if s.Items[i].HasPermission(req) {
			return true, s.Items[i]
		}
	}

	return false, nil
}

// 单个角色如何判断有没有权限
func (r *Role) HasPermission(req *PermissionRequest) bool {
	for i := range r.Spec.Permissions {
		if r.Spec.Permissions[i].HasPermission(req) {
			return true
		}
	}
	return false
}

// 判断角色里面的Permission是不是有权限
func (p *Permission) HasPermission(req *PermissionRequest) bool {
	// 确认是不是同一个服务
	if p.Service != req.Service {
		return false
	}

	// 放行所有功能
	if p.AllowAll {
		return true
	}

	// 序号查询所有的功能,确认是否允许
	for i := range p.Featrues {
		f := p.Featrues[i]
		if f.Resource == req.Resource && f.Action == req.Action {
			return true
		}
	}

	return false
}

func NewRoleSet() *RoleSet {
	return &RoleSet{
		Items: []*Role{},
	}
}

func (s *RoleSet) Add(item *Role) {
	s.Items = append(s.Items, item)
}

func NewPermissionRequest(service, resource, action string) *PermissionRequest {
	return &PermissionRequest{
		Service:  service,
		Resource: resource,
		Action:   action,
	}
}

func (req *CreateRoleRequest) Validate() error {
	return validate.Struct(req)
}

func NewRole(req *CreateRoleRequest) *Role {
	return &Role{
		Id:       xid.New().String(), //uuid
		CreateAt: time.Now().UnixMilli(),
		Spec:     req,
	}
}

func NewCreateRoleRequest() *CreateRoleRequest {
	return &CreateRoleRequest{
		Permissions: []*Permission{},
		Meta:        map[string]string{},
	}
}

func NewDefaultRole() *Role {
	return &Role{
		Spec: NewCreateRoleRequest(),
	}
}
