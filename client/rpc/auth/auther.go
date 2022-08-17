package auth

import (
	"github.com/zhou-lincong/keyauth/apps/audit"
	"github.com/zhou-lincong/keyauth/apps/policy"
	"github.com/zhou-lincong/keyauth/apps/token"
	"github.com/zhou-lincong/keyauth/client/rpc"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

func NewKeyauthAuther1(auth token.ServiceClient) *KeyauthAuther {
	return &KeyauthAuther{
		auth: auth,
		log:  zap.L().Named("http.auther"),
	}
}

func NewKeyauthAuther(client *rpc.ClientSet, serviceName string) *KeyauthAuther {
	return &KeyauthAuther{
		auth:        client.Token(),
		log:         zap.L().Named("http.auther"),
		perm:        client.Policy(),
		serviceName: serviceName,
		audit:       client.Audit(),
	}
}

// 由Keyauth提供的 HTTP认证中间件
type KeyauthAuther struct {
	log         logger.Logger
	auth        token.ServiceClient
	perm        policy.RPCClient
	serviceName string
	audit       audit.RPCClient
}
