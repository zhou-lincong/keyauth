package rpc_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/zhou-lincong/keyauth/apps/audit"
	"github.com/zhou-lincong/keyauth/apps/policy"
	"github.com/zhou-lincong/keyauth/apps/token"
	"github.com/zhou-lincong/keyauth/client/rpc"

	mcenter "github.com/infraboard/mcenter/client/rpc"
	"github.com/stretchr/testify/assert"
)

// keyauth 客户端
// 需要配置注册中心的地址
// 然后获取注册中心的客户端，使用注册中心的客户端 查询 keyauth的地址
func TestBookQuery(t *testing.T) {
	should := assert.New(t)

	conf := mcenter.NewDefaultConfig()
	conf.Address = os.Getenv("MCENTER_ADDRESS")
	conf.ClientID = os.Getenv("MCENTER_CDMB_CLINET_ID")
	conf.ClientSecret = os.Getenv("MCENTER_CMDB_CLIENT_SECRET")

	// 设置GRPC服务地址
	// conf.SetAddress("127.0.0.1:8050")
	// 携带认证信息
	// conf.SetClientCredentials("secret_id", "secret_key")

	// 传递Mcenter配置, 客户端通过Mcenter进行搜索, New一个用户中心的客户端
	// c, err := rpc.NewClient(conf)
	// if should.NoError(err) {
	// 	resp, err := c.Book().QueryBook(
	// 		context.Background(),
	// 		book.NewQueryBookRequest(),
	// 	)
	// 	should.NoError(err)
	// 	fmt.Println(resp.Items)
	// 要在运行keyauth的时候才能测试
	// === RUN   TestBookQuery
	// 2022-07-14T20:17:03.375+0800	INFO	[mcenter resolver]	resolver/resolver.go:120	search application address: [{
	// 	"Addr": "127.0.0.1:18050",
	// 	"ServerName": "",
	// 	"Attributes": null,
	// 	"BalancerAttributes": null,
	// 	"Type": 0,
	// 	"Metadata": null
	//   }]
	//   []
	//   --- PASS: TestBookQuery (0.04s)
	//   PASS
	//   ok  	github.com/zhou-lincong/keyauth/client/rpc	1.282s
	// }

	// 传递Mcenter配置, 客户端通过Mcenter进行搜索, New一个用户中心的客户端
	keyauthClient, err := rpc.NewClient(conf)

	// 使用SDK 调用Keyauth进行 凭证的校验
	// c.Token().ValidateToken()

	// 进行服务功能注册
	// keyauthClient.Endpoint().RegistryEndpoint()

	// 鉴权校验
	// keyauthClient.Policy().ValidatePermission()

	if should.NoError(err) {
		resp, err := keyauthClient.Token().ValidateToken(
			context.Background(),
			token.NewValidateTokenRequest("yTGTAj3fnPWqXIEkuicr57bf1"),
		)
		should.NoError(err)
		fmt.Println(resp)
	}
}

// 测试鉴权
//测试这个之前先post policy 传username/role
func TestValidatePermission(t *testing.T) {
	should := assert.New(t)

	conf := mcenter.NewDefaultConfig()
	conf.Address = os.Getenv("MCENTER_ADDRESS")
	conf.ClientID = os.Getenv("MCENTER_CDMB_CLINET_ID")
	conf.ClientSecret = os.Getenv("MCENTER_CMDB_CLIENT_SECRET")

	// 传递Mcenter配置, 客户端通过Mcenter进行搜索, New一个用户中心的客户端
	keyauthClient, err := rpc.NewClient(conf)
	if should.NoError(err) {
		req := policy.NewValidatePermissionRequest()
		req.Username = "member" //成员
		req.Service = "CMDB"
		req.Resource = "secret"
		req.Action = "delete" //delete

		p, err := keyauthClient.Policy().ValidatePermission(context.TODO(), req)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(p)
	}
	// === RUN   TestValidatePermission
	// 2022-07-29T19:16:37.354+0800	INFO	[mcenter resolver]	resolver/resolver.go:120	search application address: [{
	//   "Addr": "127.0.0.1:18050",
	//   "ServerName": "",
	//   "Attributes": null,
	//   "BalancerAttributes": null,
	//   "Type": 0,
	//   "Metadata": null
	// } {
	//   "Addr": "127.0.0.1:18050",
	//   "ServerName": "",
	//   "Attributes": null,
	//   "BalancerAttributes": null,
	//   "Type": 0,
	//   "Metadata": null
	// } {
	//   "Addr": "127.0.0.1:18050",
	//   "ServerName": "",
	//   "Attributes": null,
	//   "BalancerAttributes": null,
	//   "Type": 0,
	//   "Metadata": null
	// } {
	//   "Addr": "127.0.0.1:18050",
	//   "ServerName": "",
	//   "Attributes": null,
	//   "BalancerAttributes": null,
	//   "Type": 0,
	//   "Metadata": null
	// }]
	// 	e:\goproject\keyauth\client\rpc\client_test.go:103: id:"cbhs2nfj8ck2tt51rql0"  create_at:1659093341733  spec:{username:"member"  role:"member"}
	// --- PASS: TestValidatePermission (0.06s)
	// PASS
	// ok  	github.com/zhou-lincong/keyauth/client/rpc	1.298s

	//将req.Action就改成delete后测试
	// e:\goproject\keyauth\client\rpc\client_test.go:101: rpc error: code = Unknown desc = not permission access service CMDB resource secret action delete
}

func TestAuditOperate(t *testing.T) {
	should := assert.New(t)

	conf := mcenter.NewDefaultConfig()
	conf.Address = os.Getenv("MCENTER_ADDRESS")
	conf.ClientID = os.Getenv("MCENTER_CDMB_CLINET_ID")
	conf.ClientSecret = os.Getenv("MCENTER_CMDB_CLIENT_SECRET")

	// 传递Mcenter配置, 客户端通过Mcenter进行搜索, New一个用户中心的客户端
	keyauthClient, err := rpc.NewClient(conf)
	if should.NoError(err) {
		req := audit.NewOperateLog("member", "secret", "delete")
		p, err := keyauthClient.Audit().AuditOperate(context.TODO(), req)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(p)
	}
}

func init() {
	// 提前加载好 mcenter客户端, resolver需要使用
	err := mcenter.LoadClientFromEnv()
	if err != nil {
		panic(err)
	}
}
