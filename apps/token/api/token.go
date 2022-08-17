package api

import (
	"github.com/zhou-lincong/keyauth/apps/token"
	"github.com/zhou-lincong/keyauth/common/utils"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/response"
)

func (h *handler) IssueToken(r *restful.Request, w *restful.Response) {
	req := token.NewIssueTokenRequest()
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	set, err := h.service.IssueToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	response.Success(w.ResponseWriter, set)
}

func (h *handler) ValidateToken(r *restful.Request, w *restful.Response) {
	// Token哪里获取?
	// 1. URL Query String ?像飞书的一些临时token都可以通过url传过来
	// 2. Custom Header ?最好不要放在body里面
	// 3. Authorization Header

	// accessToken := ""
	// // Authorization Header
	// auth := r.HeaderParameter("Authorization")
	// // 需要把Bearer后面的切开
	// al := strings.Split(auth, " ")
	// if len(al) > 1 {
	// 	accessToken = al[1]
	// } else {
	// 	// 兼容 Authorization 后面直接跟<token>的写法
	// 	accessToken = auth
	// }
	accessToken := utils.GetToken(r.Request)

	req := token.NewValidateTokenRequest(accessToken)
	ins, err := h.service.ValidateToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}

	response.Success(w.ResponseWriter, ins)
}

func (h *handler) RevolkToken(r *restful.Request, w *restful.Response) {
	req := token.NewRevolkTokenRequest()
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	set, err := h.service.RevolkToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w.ResponseWriter, err)
		return
	}
	response.Success(w.ResponseWriter, set)
}
