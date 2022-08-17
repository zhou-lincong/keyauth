package impl

import (
	"context"

	"github.com/zhou-lincong/keyauth/apps/user"
	"github.com/zhou-lincong/keyauth/common/utils"

	"github.com/infraboard/mcube/exception"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *impl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate create user error, %s", err)
	}

	ins := user.NewUser(req)

	ins.Data.Password = utils.HashPassword(ins.Data.Password)

	// s.col.InsertMany()
	if _, err := i.col.InsertOne(ctx, ins); err != nil {
		return nil, exception.NewInternalServerError("inserted user(%s) document error, %s",
			ins.Data.Name, err)
	}
	return ins, nil
}

func (i *impl) QueryUser(ctx context.Context, req *user.QueryUserRequest) (*user.UserSet, error) {
	query := newQueryRequest(req)
	return i.query(ctx, query)
}

func (i *impl) DescribeUser(ctx context.Context, req *user.DescribeUserRequest) (*user.User, error) {
	return i.get(ctx, req)
}

func (i *impl) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}

func (i *impl) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
