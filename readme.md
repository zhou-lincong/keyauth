# keyauth
用户中心

## 架构图

## 项目说明
```
├── protocol                       # 脚手架功能: rpc / http 功能加载
│   ├── grpc.go              
│   └── http.go    
├── client                         # 脚手架功能: grpc 客户端实现 
│   ├── client.go              
│   └── config.go    
├── cmd                            # 脚手架功能: 处理程序启停参数，加载系统配置文件
│   ├── root.go             
│   └── start.go                
├── conf                           # 脚手架功能: 配置文件加载
│   ├── config.go                  # 配置文件定义
│   ├── load.go                    # 不同的配置加载方式
│   └── log.go                     # 日志配置文件
├── dist                           # 脚手架功能: 构建产物
├── etc                            # 配置文件
│   ├── xxx.env
│   └── xxx.toml
├── apps                            # 具体业务场景的领域包
│   ├── all
│   │   |-- grpc.go                # 注册所有GRPC服务模块, 暴露给框架GRPC服务器加载, 注意 导入有先后顺序。  
│   │   |-- http.go                # 注册所有HTTP服务模块, 暴露给框架HTTP服务器加载。                    
│   │   └── internal.go            #  注册所有内部服务模块, 无须对外暴露的服务, 用于内部依赖。 
│   ├── book                       # 具体业务场景领域服务 book
│   │   ├── http                   # http 
│   │   │    ├── book.go           # book 服务的http方法实现，请求参数处理、权限处理、数据响应等 
│   │   │    └── http.go           # 领域模块内的 http 路由处理，向系统层注册http服务
│   │   ├── impl                   # rpc
│   │   │    ├── book.go          # book 服务的rpc方法实现，请求参数处理、权限处理、数据响应等 
│   │   │    └── impl.go           # 领域模块内的 rpc 服务注册 ，向系统层注册rpc服务
│   │   ├──  pb                    # protobuf 定义
│   │   │     └── book.proto       # book proto 定义文件
│   │   ├── app.go                 # book app 只定义扩展
│   │   ├── book.pb.go             # protobuf 生成的文件
│   │   └── book_grpc.pb.go        # pb/book.proto 生成方法定义
├── version                        # 程序版本信息
│   └── version.go                    
├── README.md                    
├── main.go                        # Go程序唯一入口
├── Makefile                       # make 命令定义
└── go.mod                         # go mod 依赖定义
```
client                              #
    rest                            # restful客户端
    rpc                             # rpc客户端
        client.go   
后面rpc会通过注册中心做微服务发现      

## 快速开发
make脚手架
```sh
➜  keyauth git:(master) ✗ make help
dep                            Get the dependencies
lint                           Lint Golang files
vet                            Run go vet
test                           Run unittests
test-coverage                  Run tests with coverage
build                          Local build
linux                          Linux build
run                            Run Server
clean                          Remove previous build
help                           Display this help screen
```
1. 使用安装依赖的Protobuf库(文件)
```sh
# 把依赖的probuf文件复制到/usr/local/include

# 创建protobuf文件目录
$ make -pv /usr/local/include/github.com/infraboard/mcube/pb

# 找到最新的mcube protobuf文件
$ ls `go env GOPATH`/pkg/mod/github.com/infraboard/

# 复制到/usr/local/include
$ cp -rf pb  /usr/local/include/github.com/infraboard/mcube/pb
```

2. 添加配置文件(默认读取位置: etc/keyauth.toml)
```sh
$ 编辑样例配置文件 etc/keyauth.toml.book
$ mv etc/keyauth.toml.book etc/keyauth.toml
```

3. 启动服务
```sh
# 编译protobuf文件, 生成代码
$ make gen
# 如果是MySQL, 执行SQL语句(docs/schema/tables.sql)
$ make init
# 下载项目的依赖
$ make dep
# 运行程序
$ make run
```

## 相关文档


## 认证
### HTTP接口的认证(OpenAPI)
该接口会被哪些用户使用:
+ Web界面 ---> OpenAPI
+ 其他第三方服务, 监控系统(Prometheus), 监控资源发现

如何认证:
+ basic auth: (user/password) ---> user:
    + 存放的Header： Authorization
    + "Basic " +  user:password 的base64编码
    + header: Authorization = "Basic base64(user:password)"
+ 存在的问题：用户的密码 在每次API调用的时候都会把明文传输, 甚至可以登录Web 去修改你的密码

系统 -->（来调用） 接口 (编程用户: 程序使用)，很多云厂商会区分成：
+ 需要登录界面的用户
+ 编程用户

最好的方式 不直接使用用户的 user/password, 可以用户颁发一个代表自己身份的令牌(token)
即使用户发现自己的令牌被盗取了, 可以吊销该令牌
怎么基于Token做一个统一认证系统? 我们总不至于每个系统 都自己保存一个user表

现在的系统目前就cmdb资源中心，后面会有workflow分布式、audit审计中心
通过web的ifree直接进行页面嵌套打通各个系统的用户系统，交互逻辑不友好，后续系统很难进行整合


### GRPC接口的认证


### 用户中心
流程：用户login输入user/password去用户中心验证，验证通过返回用户身份的token。
    然后用户携带这个token去访问cmdb，cmdb通过grpc接口去用户中心确认token的合法性。
适配认证方式：user/password、LDAP、第三方认证Oath2:钉钉/企业微信/飞书、编程用户提供方式：AK/SK、
    给编程用户提供的token。
用户表：
    username：admin
    password：123456(hash)。通过比对hash，无需比对明文密码
    domian:Org/公司/租户/域，本domian里面的资源只有本domian能访问
    namespace:命名空间，用于做资源隔离（资源+人+权限）
        > 访问范围
用户认证token：
    accessToken:随机字符串,本身不包含任何信息，称为raw token
    username：该token代表哪个用户的身份
    issueAt：颁发时间、token的使用时长限制、计算出当前token的使用时长
    accessTokenExpiredAt:失效时间（比如1个小时后），失效了怎么办？
    refreshToken：该token允许刷新一个新的token或者延长token的过期时间
    refreshTokenExpiredAt：失效时间（比如1个小时后），如果refreshToken隔了很久没用，那也有失效时间
        通过这样的设计，如果用户中间一个小时没访问，那么会话就断了。这样能知道会话时长2H。
服务表：clientId、clientSecret
服务认证：对比clientId、clientSecret
token里面的信息：

## 流程
### 项目初始化
mcube project init github.com/zhou-lincong/keyauth
不接入权限中心
mongoDb
mongoDb服务地址
认证数据库名称：keyauth（账号在哪个数据库去认证，每一个库都有认证的用户名和密码，例如有5个库db1/db2/.../db5,每个库都有用户管理表，还有一个admin用户表，在admin空间下面。这个时候用admin去验证就能操作所有db，所以专门建一个keyauth用户）
认证用户：keyauth
数据库名称：keyauth
生成样例代码
选择gin-restful
make dep
make pb(路径故障根据提示手动解决)

### mongoDB初始化
1. 采用docker安装
docker pull mongo
docker run -itd -p 27017:27017 mongo  // --name mongodb
docker ps  //获取CONTAINER ID
docker exec -it e4eec90f90af bash
mongo
    2. 编辑 /etc/mongod.conf 开启认证访问(可选), 开启后要重启下服务
    security:
    authorization: enabled
3. 创建管理员账号
use admin
db.createUser({user:"admin",pwd:"123456",roles:["root"]})
db.auth("admin", "123456")
4. 添加库用户
use keyauth
db.createUser({user: "keyauth", pwd: "123456", roles: [{ role: "dbOwner", db: "keyauth" }]})
db.auth("keyauth", "123456")
5. 退出exit

创建mcenter
1. docker ps
2. docker exec -it e4eec90f90af mongo
3. use mcenter
4. db.createUser({user: "mcenter", pwd: "123456", roles: [{ role: "dbOwner", db: "mcenter" }]})
5. db.auth("mcenter", "123456")

## 逻辑梳理
除了cmdb，后面还有workflow分布式流水线系统、audit审计中心、用户中心（权限系统）。
    cmdb单体应用做完后，就需要做成一个分布式架构。如果这个用户中心放在cmdb，那么就在apps里面增加一个token或user的模块，由它来实现认证的功能。但这样的话，这个用户中心就没办法扩展。如果以后要增加工单系统、审计系统，这样就不可能全加到cmdb里面去或者去cmdb里面去做认证，也不可能在这两个系统里面再增加用户系统。所以依照微服务的思想，把公共的功能抽离出来，成为独立的服务，通过微服务调用的方式，让它们对接上。比如审计中心，因为可能每一个接口都需要审计。用户中心，所有的模块都到用户中心进行认证。（cmdb到用户中心的转变，接下来做分布式。
    cmdb单体应用做完后，就需要做成一个分布式架构。如果这个用户中心放在cmdb，那么就在apps里面增加一个token或user的模块，由它来实现认证的功能。但这样的话，这个用户中心就没办法扩展。如果以后要增加工单系统、审计系统，这样就不可能全加到cmdb里面去或者去cmdb里面去做认证，也不可能在这两个系统里面再增加用户系统。所以依照微服务的思想，把公共的功能抽离出来，成为独立的服务，通过微服务调用的方式，让它们对接上。比如审计中心，因为可能每一个接口都需要审计。用户中心，所有的模块都到用户中心进行认证。（cmdb到用户中心的转变，接下来做分布式。）
最终的界面：是多个模块出完了之后，统一走一个代理过后做访问。
    比如前端页面，有cmdb系统、有审计中心、用户中心、workflow流水线相关的模块的一个页面。人后每一个页面都会访问后端不同的系统。不同的系统之间，比如用户要访问cmdb的资源，没有认证系统，就是无论是谁通过api都能调用把数据拿到。如果要做认证，就不能让cmdb自己做认证。
    第一步，前端需要携带一个token，也就是用户的凭证，这个token可能是通过password携带，每个后端的系统认的是token，token代表的是用户的身份。除了可以通过user、password这种方式获得token，还可以用其他的方式获得，比如LDAP、Oath第三方认证、AK/SK、临时token，然后去访问后端的系统。比如去访问cmdb，但cmdb不知道这个用户是谁或合不合法，就需要拿着这个token去用户中心验证，因此用户中心需要实现一个validata_token的方法，然后让其他模块通过rpc来调用这个方法。调用完了，验证token是合法的，就把请求返回给用户。那在前面拦截token的时候，同时也要实现HTTP层的中间件。中间件的逻辑就是拦截token，调用rpc的服务验证。通过restful API接口返回代表用户信息的token（验证完的结果），把用户的信息放在请求的上下文当中，交给cmdb去访问。
实现流程：1.用户的认证：
    用户中心需要实现用户认证的接口，需要提供简单的用户名密码，然后返回一个token出来。这个token就可以访问资源cmdb，访问过后就有中间件拿着token调用rpc接口去用户中心验证合不合法，验证完之后返回这个token代表的用户是谁有什么权限配置等，这样就可以把权限的逻辑独立开了。内部走的rpc，外部与web走的HTTP。
          2.服务的认证：
    确认调用方的ID和secret的正确性、由谁来颁发凭证？
    外部的认证走用户中心，内部的认证走注册中心做一个简单的认证。
    注册中心逻辑：用clientID和clientsecret做验证，需要创建服务。
        创建1.cmdb、2.keyauth，这两个服务。然后为这两个服务颁发GRPC客户端凭证。
        当cmdb对用户中心发起访问的时候，就需要做一个鉴权，由用户中心去注册中心确认客户端凭证是否正确，OK了之后，cmdb就能调用户中心了。
        注册中心只管理这用户的凭证，不管理用户的实例。
    注册中心和用户中心的功能是很接近的，用户中心是用用户名和密码换token、校验的是用户和业务逻辑，注册中心是通过clientID和secret检验服务的逻辑是否OK的、检验内部服务调用如A服务是否能调用B服务，
架构一个系统时，需要考虑其中的哪个服务应不应该拆分开，拆分开了之后，谁应该提供sdk，内部调用方是谁，调用的权限问题，

## Mongodb
1. 如何连接MongoDB, 以及MongoDB的配置
2. 基于MongoDB的 CRUD
    + bson struct tag: 用于MongoDB 保存时，完成Object -> bson的映射
    + "_id": 内置的BSON TAG, 类似于MySQL主键, _: 代表的倒序索引, 不用额外多创建一个Index, 最好沿用 
    + 通过DB对象的Collection对象来 进行CRUD操作
    + 

## 对接注册中心
注册中心使用的mcenter服务: https://github.com/infraboard/mcenter

主题: 服务发现
 + 服务注册
    1. 添加注册中心的配置
    2. 然后初始化全局的 注册中心客户端实例 rpc.C()
    3. GRPC 服务启动时 调用 注册中心的客户端把 当前GPRC监听的地址注册过去
    4. 当GRPC服务Stop时, 注销注册中心的实例
 + 服务解索(GRPC Client)
    1. 通过GRPC 的NamedResolver 来进行服务的发现
    2. 也是加载 注册中心 的 GRPC客户端: rpc.C() --? , 
        因为 Mcenter 提供的Resolver需要依赖 注册中心的客户端来进行 服务实例的搜索
    3. 在服务启动的时候 初始化的时候 就完成以上步骤
    4. GRPC客户端 配置注册中心的访问凭证, 以及需要访问的服务的名称, Resovler就能完成服务名称--》 
       地址的解析
namingClient：名称解析的客户端，通过RegisterInstance把IP、port、服务的名称注册过去，最后就可以在GetService通过服务的名称serviceName找到注册在里面的地址，从而不必每写一个服务都要配置一下地址。不然这样就会导致在服务调用的关系复杂且数量多的时候，每个服务都去配置就会很麻烦，而且sdk也难写。
为了把调用关系解耦，跟IOC的逻辑类似，剥离出一个service registry统一的管理的地方，管理所有服务的地址当需要哪个地址的时候，就去这个第三方获取，也就是这个service注册中心，然后带着地址去访问。
还包含其他的服务，例如服务的重试、限流、负载均衡、健康检查、配置、调用权限。
注册后也可以注销。
在服务启动能正常提供服务后，再把服务注册给注册中心，避免出现有些服务调度过来，不能提供服务的请求异常。
GRPC怎么对接注册中心：
> 正常情况下，grpc会有builder和resolver两个概念。
builder，就是为了构造一个resolver对象，然后这个对象会去注册中心搜索，根据相应的条件做地址解析。解析完后，会返回具体的地址，比如传了一个name:keyauth的地址过去。第二步，拿到地址后就建立连接，grpc客户端就能访问grpc服务端。之前地址是手动写死的，现在有了注册中心，地址就不是手动写死的，通过注册中心做自动发现。那grpc客户端里面就需要一种机制，自动发现地址然后通过地址去连接。
> 演示：demo客户端写了一个mcenter resolver，由这个mcenter resolver去向注册中心发起一次地址解析，解析的时候带上服务的名字。如解析keyauth服务的地址，那么它会根据当前数据库里面，有没有keyauth的地址，找到后返给客户端。这里的地址可能是一条，可能是多条，这个时候grpc的server会从中挑选几个，挑选规格有很多有RR、PF挑选第一个，可以server端进行排序。实现注册中心，会方便服务端或者客户端对返回的结果进行排序。用开源的注册中心，可能不好控制，如根据权重排序。


准备完成 Keyauth的客户端, CMDB 就可以通过该客户端来和Keyauth进行交互
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

	// 把Mcenter的配置传递给Keyauth的客户端
	c, err := rpc.NewClient(conf)

	// 使用SDK 调用Keyauth进行 凭证的校验
	// c.Token().ValidateToken()

	if should.NoError(err) {
		resp, err := c.Token().ValidateToken(
			context.Background(),
			token.NewValidateTokenRequest("yTGTAj3fnPWqXIEkuicr57bf1"),
		)
		should.NoError(err)
		fmt.Println(resp)
	}
}
```
cmdb 如何使用Keyauth的客户端进行 Token的校验, 需要一个HTTP 的 认证中间件:
因为Keyauth是用户中心, HTTP的权限中间件需要和keyauth交互, 依赖Keyatuh的SDK, 因此这个中间件由 keyauth提供，在client/rpc/增加auth/server.go


#### 代码实现
目前已经完成用户表和token表，校验逻辑也完成了
接下来需要实现认证中间件，让cmdb那边加上一个中间件，让它所有的token访问cmdb的时候都必须去用户中心来认证一遍。
需要用rpc的客户端去做一个对接，由cmdb服务通过grpc去用户中心keyauth进行token服务的校验，但并不是谁都能来校验，因此涉及到注册中心。
注册中心前期管理是服务如用户中心、cmdb，这些服务的clientID和clientsecret，是在管理这些服务的凭证。跟用户中心是一样的，凭证拿到注册中心来验证，合法了然后再让用户调具体的用户中心服务。因为注册中心只是验证了一下client凭证，并没有做过多的逻辑，后期也可以扩展如是否允许cmdb调用户中心、调workflow，都可以再内部做控制。
>主要逻辑：提供grpc客户端，把服务的实例注册给注册中心，再由其他服务访问注册中心获取服务实例的地址
