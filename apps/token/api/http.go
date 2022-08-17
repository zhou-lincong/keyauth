package api

import (
	"github.com/zhou-lincong/keyauth/apps/token"

	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	h = &handler{}
)

type handler struct {
	service token.ServiceServer
	log     logger.Logger
}

func (h *handler) Config() error {
	h.log = zap.L().Named(token.AppName)
	h.service = app.GetGrpcApp(token.AppName).(token.ServiceServer)
	return nil
}

func (h *handler) Name() string {
	return token.AppName
}

func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry(ws *restful.WebService) {
	tags := []string{h.Name()}

	ws.Route(ws.POST("/issue").To(h.IssueToken).
		Doc("issue token").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(token.IssueTokenRequest{}).
		Writes(response.NewData(token.Token{})))

	ws.Route(ws.GET("/validate").To(h.ValidateToken).
		Doc("validate token").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata("action", "get"))

	ws.Route(ws.POST("/revolk").To(h.RevolkToken).
		Doc("revolk token").
		Param(ws.PathParameter("id", "identifier of the user").DataType("integer").DefaultValue("1")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Metadata("action", "delete"))

}

func init() {
	app.RegistryRESTfulApp(h)
}
