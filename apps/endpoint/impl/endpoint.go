package impl

import (
	"context"

	"github.com/zhou-lincong/keyauth/apps/endpoint"
)

// 实现服务注册接口的具体逻辑
func (s *service) RegistryEndpoint(ctx context.Context, req *endpoint.EndpointSet) (
	*endpoint.RegistryResponse, error) {
	if err := s.save(ctx, req); err != nil {
		return nil, err
	}
	resp := endpoint.NewRegistryResponse()
	return resp, nil
}
