package rpc

import (
	"fmt"

	"github.com/zhou-lincong/keyauth/apps/audit"
	"github.com/zhou-lincong/keyauth/apps/endpoint"
	"github.com/zhou-lincong/keyauth/apps/policy"
	"github.com/zhou-lincong/keyauth/apps/role"
	"github.com/zhou-lincong/keyauth/apps/token"

	"github.com/infraboard/mcenter/client/rpc"
	"github.com/infraboard/mcenter/client/rpc/auth"
	"github.com/infraboard/mcenter/client/rpc/resolver"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client *ClientSet
)

// SetGlobal todo
func SetGlobal(cli *ClientSet) {
	client = cli
}

// C Global
func C() *ClientSet {
	return client
}

// NewClient todo 传注册中心的地址过去
func NewClient(conf *rpc.Config) (*ClientSet, error) {
	zap.DevelopmentSetup()
	log := zap.L()
	// resolver 进行解析的时候 需要mcenter客户端实例已经初始化
	conn, err := grpc.Dial(
		// 先把url逻辑补上  127.0.0.1:18010 GRPC server端的地址
		// 基于服务发现  Dial to "passthrough://  dns://keyauth.org "mcenter://keyauth",
		// resolver会把这些地址带上实例的ClientID和ClientSecret去访问注册中心
		//从而拿到注册中心里面的地址，进行交互
		//但是mcenter的服务还没起来和初始化，后面再处理
		fmt.Sprintf("%s://%s", resolver.Scheme, "keyauth"), //代替 conf.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 代替 grpc.WithPerRPCCredentials(conf.Authentication),auth.NewAuthentication
		grpc.WithPerRPCCredentials(auth.NewAuthentication(conf.ClientID, conf.ClientSecret)),
		// 直到连接成功才返回
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &ClientSet{
		conn: conn,
		log:  log,
	}, nil
}

// Client 客户端
type ClientSet struct {
	conn *grpc.ClientConn
	log  logger.Logger
}

// Book服务的SDK
// func (c *ClientSet) Book() book.ServiceClient {
// 	return book.NewServiceClient(c.conn)
// }

// token服务的SDK
func (c *ClientSet) Token() token.ServiceClient {
	return token.NewServiceClient(c.conn)
}

// Endpoint服务的SDK
func (c *ClientSet) Endpoint() endpoint.RPCClient {
	return endpoint.NewRPCClient(c.conn)
}

// role服务的SDK
func (c *ClientSet) Role() role.RPCClient {
	return role.NewRPCClient(c.conn)
}

// policy服务的SDK
func (c *ClientSet) Policy() policy.RPCClient {
	return policy.NewRPCClient(c.conn)
}

// Audit服务的SDK
func (c *ClientSet) Audit() audit.RPCClient {
	return audit.NewRPCClient(c.conn)
}
