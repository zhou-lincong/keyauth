package api

import (
	"github.com/zhou-lincong/keyauth/apps/policy"

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
	service policy.Service
	log     logger.Logger
}

func (h *handler) Config() error {
	h.log = zap.L().Named(policy.AppName)
	h.service = app.GetGrpcApp(policy.AppName).(policy.Service)
	return nil
}

func (h *handler) Name() string {
	return policy.AppName
}

func (h *handler) Version() string {
	return "v1"
}

func (h *handler) Registry(ws *restful.WebService) {
	tags := []string{h.Name()}

	ws.Route(ws.POST("/").To(h.CreatePolicy).
		Doc("create a policy").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(policy.CreatePolicyRequest{}).
		Writes(response.NewData(policy.Policy{})))
}

func init() {
	app.RegistryRESTfulApp(h)
}
