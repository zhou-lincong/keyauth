# Keyauth SDK

```go
// keyauth 客户端
// 需要配置注册中心的地址
// 获取注册中心的客户端，使用注册中心的客户端 查询 keyauth的地址
func TestBookQuery(t *testing.T) {
	should := assert.New(t)

	conf := mcenter.NewDefaultConfig()
	conf.Address = os.Getenv("MCENTER_ADDRESS")
	conf.ClientID = os.Getenv("MCENTER_CDMB_CLINET_ID")
	conf.ClientSecret = os.Getenv("MCENTER_CMDB_CLIENT_SECRET")

	// 传递Mcenter配置, 客户端通过Mcenter进行搜索, New一个用户中心的客户端
	keyauthClient, err := rpc.NewClient(conf)

	// 使用SDK 调用Keyauth进行 凭证的校验
	// c.Token().ValidateToken()

	if should.NoError(err) {
		resp, err := keyauthClient.Token().ValidateToken(
			context.Background(),
			token.NewValidateTokenRequest("yTGTAj3fnPWqXIEkuicr57bf1"),
		)
		should.NoError(err)
		fmt.Println(resp)
	}
}
```

客户端是有版本概念, 通过
```
# 当前客户端的版本 v0.0.1
git tag v0.0.1

# 其他服务使用时, go mod里面指定版本即可
```

使用 go restful 客户端编写SDK
```go
i.client.Post("/sdsfd").
	Header("token", "xxx").
	Param("page_size", 1).
	Timeout(3 * time.Second).
	Body(nil).
	Do(context.TODO()).
	Into(resp).Error()
```