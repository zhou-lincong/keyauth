package auth

import (
	"time"

	"github.com/zhou-lincong/keyauth/apps/audit"
	"github.com/zhou-lincong/keyauth/apps/policy"
	"github.com/zhou-lincong/keyauth/apps/token"
	"github.com/zhou-lincong/keyauth/common/utils"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
)

// Go-Restful auth Middleware
// 通过 r.Filter(<middleware>) 添加中间件
func (a *KeyauthAuther) RestfulAuthHandlerFunc(
	req *restful.Request,
	resp *restful.Response,
	chain *restful.FilterChain,
) {
	//1.把路由里面的装饰信息获取出来
	meta := req.SelectedRoute().Metadata()
	a.log.Debug(meta)

	//获取meta信息,可以为这个meta封装一个对象，实现一个get方法，让代码看起来优雅一点
	//判断是否开启认证
	var isAuthEnable bool
	if authV, ok := meta[label.Auth]; ok {
		switch v := authV.(type) {
		case bool:
			isAuthEnable = v
		case string:
			isAuthEnable = v == "true"
		}
	}

	if isAuthEnable {
		// 处理Request对象
		// *restful.Request

		// 1. 认证中间件, 需要获取Token
		tkStr := utils.GetToken(req.Request)

		// 2. 到用户中心验证token合法性, 依赖用户中心的Grpc Client
		validateReq := token.NewValidateTokenRequest(tkStr)
		tk, err := a.auth.ValidateToken(req.Request.Context(), validateReq)
		if err != nil {
			response.Failed(resp.ResponseWriter, exception.NewUnauthorized(err.Error()))
			return
		}

		// 处理完成后 需要直接中断返回
		// *restful.Response
		// resp.Header()
		// resp.Write()

		// 如果处理ok需要把一些中间结果 然后后面的请求也能方式, 需要把结果保存中 上下文中
		// 原生 c.Request.Context()
		// 如何:  context Get
		// 认证完成后, 我们需要用户名称或者其他信息 传递下去

		// Set Context
		req.SetAttribute("token", tk)
		// req.Attribute("token").(*token.Token)
	}

	// 获取meta信息, get, 判断是否开启鉴权
	var isPermEnable bool
	if authV, ok := meta[label.Permission]; ok {
		switch v := authV.(type) {
		case bool:
			isPermEnable = v
		case string:
			isPermEnable = v == "true"
		}
	}

	// 认证后，才能鉴权
	if isAuthEnable && isPermEnable {
		//从req里面获取token
		tk := req.Attribute("token").(*token.Token)
		permReq := policy.NewValidatePermissionRequest()
		permReq.Service = a.serviceName
		permReq.Username = tk.Data.Username

		if meta != nil {
			if v, ok := meta[label.Resource]; ok {
				permReq.Resource, _ = v.(string)
			}
			if v, ok := meta[label.Action]; ok {
				permReq.Action, _ = v.(string)
			}
		}

		_, err := a.perm.ValidatePermission(req.Request.Context(), permReq)
		if err != nil {
			response.Failed(resp.ResponseWriter, exception.NewPermissionDeny(err.Error()))
			return
		}
	}

	// 获取meta信息, get, 判断是否开启鉴权
	var isAuditEnable bool
	if authV, ok := meta[label.Audit]; ok {
		switch v := authV.(type) {
		case bool:
			isAuditEnable = v
		case string:
			isAuditEnable = v == "true"
		}
	}

	start := time.Now()

	// chain 用于将请求flow下去
	chain.ProcessFilter(req, resp)

	//花费的时间
	cost := time.Now().Sub(start).Milliseconds()

	// 如果有记录Respone需求的话，需要在Process才有进行的
	// 认证后，才能审计
	if isAuthEnable && isAuditEnable {
		tk := req.Attribute("token").(*token.Token)
		auditReq := audit.NewOperateLog(tk.Data.Username, "", "")
		auditReq.Service = a.serviceName
		auditReq.Url = req.Request.URL.String()
		auditReq.Cost = cost
		auditReq.StatusCode = int64(resp.StatusCode())
		auditReq.UserAgent = req.Request.UserAgent()
		// X-Forwar-For
		// auditReq.RemoteIp = req.Request.RemoteAddr
		auditReq.RemoteIp = request.GetRemoteIP(req.Request)

		if meta != nil {
			if v, ok := meta[label.Resource]; ok {
				auditReq.Resource, _ = v.(string)
			}
			if v, ok := meta[label.Action]; ok {
				auditReq.Action, _ = v.(string)
			}
		}

		_, err := a.audit.AuditOperate(req.Request.Context(), auditReq)
		if err != nil {
			a.log.Warnf("audit operate failed, %s", err)
			return
		}
	}
}
