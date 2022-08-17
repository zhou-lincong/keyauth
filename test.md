# 测试获取用户名，没有的时候返回这样的结果
GET http://127.0.0.1:8050/keyauth/api/v1/user
{
    "code": 0,
    "data": {
        "total": 0,
        "items": []
    }
}


# 测试传空的进去
POST http://127.0.0.1:8050/keyauth/api/v1/user
{
    "code": 400,
    "namespace": "global",
    "reason": "请求不合法",
    "message": "validate create user error, Key: 'CreateUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"
}


# 测试创建用户
POST http://127.0.0.1:8050/keyauth/api/v1/user
{
   "name": "袁鑫鑫",
   "password": "123456"
}
{
    "code": 0,
    "data": {
        "id": "cb78q6vj8ck0p30h2dmg",
        "create_at": 1657703707001,
        "update_at": 0,
        "update_by": "",
        "data": {
            "create_by": "",
            "name": "袁鑫鑫",
            "password": "$2a$14$Mf7JajIdy2kw1L86mmM/dOhRyNeMFq0vRdNfh3FTCx7GUNK933f/2",
            "domain": "default"
        }
    }
}


# 测试创建相同用户
{
    "code": 500,
    "namespace": "global",
    "reason": "系统内部错误",
    "message": "inserted user(袁鑫鑫) document error, write exception: write errors: [E11000 duplicate key error collection: keyauth.user index: data.domain_-1_data.name_-1 dup key: { data.domain: \"default\", data.name: \"袁鑫鑫\" }]"
}


# 测试颁发token
post http://127.0.0.1:8050/keyauth/api/v1/token
{
   "user_name": "袁鑫鑫",
   "password": "123456"
}
{
    "code": 0,
    "data": {
        "access_token": "Xm0X4HyuMp3lx9dEcPraJKif",
        "issue_at": 1657720103279,
        "update_at": 0,
        "update_by": "",
        "data": {
            "grante_type": "PASSWORD",
            "user_domain": "default",
            "user_name": "袁鑫鑫",
            "password": "123456"
        },
        "access_token_expired_at": 1657720703279,
        "refresh_token": "rHx8fVOHOvOl8I8TzV8Z5RPI6EsO9yiE",
        "refresh_token_expired_at": 1657723103279,
        "domian": "",
        "meta": null
    }
}


# 测试输入不存在的账号密码颁发token
post http://127.0.0.1:8050/keyauth/api/v1/token
{
   "user_name": "袁鑫鑫1",
   "password": "123456"
}
{
    "code": 401,
    "namespace": "global",
    "reason": "认证失败",
    "message": "user or password not correct"
}


# 测试post http://127.0.0.1:8050/keyauth/api/v1/token/issue
{
   "user_name": "袁鑫鑫",
   "password": "123456"
}
{
    "code": 0,
    "data": {
        "access_token": "FArcPEAV3oj222rQGSY7Pw9n",
        "issue_at": 1657768618594,
        "update_at": 0,
        "update_by": "",
        "data": {
            "grante_type": "PASSWORD",
            "user_domain": "default",
            "user_name": "袁鑫鑫",
            "password": "123456"
        },
        "access_token_expired_at": 1657769218594,
        "refresh_token": "kdTzQqJ9YHf7f07GRKatsv501Wt7Lw3r",
        "refresh_token_expired_at": 1657771618594,
        "domian": "",
        "meta": null
    }
}


# 测试validate token 
get http://127.0.0.1:8050/keyauth/api/v1/token/validate
authorization type bearertoken token:
{
    "code": 0,
    "data": {
        "access_token": "FArcPEAV3oj222rQGSY7Pw9n",
        "issue_at": 1657768618594,
        "update_at": 0,
        "update_by": "",
        "data": {
            "grante_type": "PASSWORD",
            "user_domain": "default",
            "user_name": "袁鑫鑫",
            "password": ""
        },
        "access_token_expired_at": 1657769218594,
        "refresh_token": "kdTzQqJ9YHf7f07GRKatsv501Wt7Lw3r",
        "refresh_token_expired_at": 1657771618594,
        "domian": "",
        "meta": null
    }
}


# 测试撤销token
POST http://127.0.0.1:8050/keyauth/api/v1/token/revolk
{
   "access_token": "FArcPEAV3oj222rQGSY7Pw9n",
   "refresh_token": "kdTzQqJ9YHf7f07GRKatsv501Wt7Lw3r"
}
{
    "code": 0,
    "data": {
        "access_token": "FArcPEAV3oj222rQGSY7Pw9n",
        "issue_at": 1657768618594,
        "update_at": 0,
        "update_by": "",
        "data": {
            "grante_type": "PASSWORD",
            "user_domain": "default",
            "user_name": "袁鑫鑫",
            "password": ""
        },
        "access_token_expired_at": 1657769218594,
        "refresh_token": "kdTzQqJ9YHf7f07GRKatsv501Wt7Lw3r",
        "refresh_token_expired_at": 1657771618594,
        "domian": "",
        "meta": null
    }
}
然后用valitate验证，在数据库找不到了


# 测试mcenter
1. 先创建本地mongo的mcenter账号
2. 去本地mcenter 执行go run main.go init 获取keyauth的client_id和client_secret，并添加到etc/unit_test.env
3. 运行本地mcenter
4. 运行keyauth
结果：
$ make run
2022-07-14T19:45:10.403+0800    INFO    [INIT]  cmd/start.go:153        log level: debug
2022-07-14T19:45:10.510+0800    INFO    [CLI]   cmd/start.go:91 loaded grpc app: [book token user]
2022-07-14T19:45:10.510+0800    INFO    [CLI]   cmd/start.go:92 loaded http app: [book token user]
2022-07-14T19:45:10.511+0800    INFO    [CLI]   cmd/start.go:94 loaded internal app: []
2022-07-14T19:45:10.511+0800    INFO    [GRPC Service]  protocol/grpc.go:72     GRPC 服务监听地址: 127.0.0.1:18050 
2022-07-14T19:45:10.528+0800    INFO    [HTTP Service]  protocol/http.go:78     Get the API using http://127.0.0.1:8050/apidocs.json
2022-07-14T19:45:10.530+0800    INFO    [HTTP Service]  protocol/http.go:81     HTTP服务启动成功, 监听地址: 127.0.0.1:8050
2022-07-14T19:45:11.527+0800    ERROR   [GRPC Service]  protocol/grpc.go:95     registry to mcenter error, rpc error: code = Unauthenticated desc = client_id or client_secret is ""
最后一行报错把参数添加到etc/config.toml重新运行：
结果：
$ make run
2022-07-14T20:04:36.280+0800    INFO    [INIT]  cmd/start.go:153        log level: debug
2022-07-14T20:04:36.386+0800    INFO    [CLI]   cmd/start.go:91 loaded grpc app: [user book token]
2022-07-14T20:04:36.387+0800    INFO    [CLI]   cmd/start.go:92 loaded http app: [book token user]
2022-07-14T20:04:36.388+0800    INFO    [CLI]   cmd/start.go:94 loaded internal app: []
2022-07-14T20:04:36.391+0800    INFO    [GRPC Service]  protocol/grpc.go:72     GRPC 服务监听地址: 127.0.0.1:18050 
2022-07-14T20:04:36.406+0800    INFO    [HTTP Service]  protocol/http.go:78     Get the API using http://127.0.0.1:8050/apidocs.json
2022-07-14T20:04:36.407+0800    INFO    [HTTP Service]  protocol/http.go:81     HTTP服务启动成功, 监听地址: 127.0.0.1:8050
2022-07-14T20:04:37.788+0800    INFO    [GRPC Service]  protocol/grpc.go:99     registry to mcenter success


# 测试创建角色
POST http://127.0.0.1:8050/keyauth/api/v1/role
{
    "name": "member",
    "description": "测试",
    "permissions": [
        {
            "service": "CMDB",
            "featrues": [
                {
                    "resource": "secret",
                    "action": "list"
                },
                {
"resource": "secret",
                    "action": "get"
                }
            ]
        }
    ]
}
返回：
{
    "code": 0,
    "data": {
        "id": "cbhqnanj8ck13c5j22j0",
        "name": 1659087786553,
        "spec": {
            "name": "member",
            "description": "测试",
            "permissions": [
                {
                    "service": "CMDB",
                    "featrues": [
                        {
                            "resource": "secret",
                            "action": "list"
                        },
                        {
                            "resource": "secret",
                            "action": "get"
                        }
                    ]
                }
            ],
            "meta": {}
        }
    }
}


# 测试创建策略
POST http://127.0.0.1:8050/keyauth/api/v1/policy
{
    "username": "member",
    "role": "member"
}


#
1. 
get  http://127.0.0.1:8060/CMDB/api/v1/secret/cb41jm7j8ck1r969o3hg
{
    "code": 0,
    "data": {
        "id": "cb41jm7j8ck1r969o3hg",
        "create_at": 1657280984492,
        "data": {
            "description": "袁鑫",
            "vendor": "TENCENT",
            "allow_regions": [
                "*"
            ],
            "crendential_type": "API_KEY",
            "address": "",
            "api_key": "woaini",
            "api_secret": "******",
            "request_rate": 5,
            "create_by": ""
        }
    }
}
2. 
get  http://127.0.0.1:8060/CMDB/api/v1/secret
{
    "code": 0,
    "data": {
        "total": 1,
        "items": [
            {
                "id": "cb41jm7j8ck1r969o3hg",
                "create_at": 1657280984492,
                "data": {
                    "description": "袁鑫",
                    "vendor": "TENCENT",
                    "allow_regions": [
                        "*"
                    ],
                    "crendential_type": "API_KEY",
                    "address": "",
                    "api_key": "woaini",
                    "api_secret": "******",
                    "request_rate": 5,
                    "create_by": ""
                }
            }
        ]
    }
}