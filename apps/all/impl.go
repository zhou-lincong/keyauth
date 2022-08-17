package all

import (
	// 注册所有GRPC服务模块, 暴露给框架GRPC服务器加载, 注意 导入有先后顺序
	_ "github.com/zhou-lincong/keyauth/apps/audit/impl"
	_ "github.com/zhou-lincong/keyauth/apps/endpoint/impl"
	_ "github.com/zhou-lincong/keyauth/apps/policy/impl"
	_ "github.com/zhou-lincong/keyauth/apps/role/impl"
	_ "github.com/zhou-lincong/keyauth/apps/token/impl"
	_ "github.com/zhou-lincong/keyauth/apps/user/impl"
)
