package api

import (
	"github.com/zhou-lincong/keyauth/apps/role"

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
	service role.Service
	log     logger.Logger
}

func (h *handler) Config() error {
	h.log = zap.L().Named(role.AppName)
	h.service = app.GetGrpcApp(role.AppName).(role.Service)
	return nil
}

func (h *handler) Name() string {
	return role.AppName
}

func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry(ws *restful.WebService) {
	tags := []string{h.Name()}

	ws.Route(ws.POST("/").To(h.CreateRole).
		Doc("create a role").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(role.CreateRoleRequest{}).
		Writes(response.NewData(role.Role{})))
}

func init() {
	app.RegistryRESTfulApp(h)
}
