package all

import (
	// 注册所有HTTP服务模块, 暴露给框架HTTP服务器加载
	_ "github.com/zhou-lincong/keyauth/apps/policy/api"
	_ "github.com/zhou-lincong/keyauth/apps/role/api"
	_ "github.com/zhou-lincong/keyauth/apps/token/api"
	_ "github.com/zhou-lincong/keyauth/apps/user/api"
)
