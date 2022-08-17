package endpoint

const (
	AppName = "endpoint"
)

type Service interface {
	RPCServer
}

//补充生成文档列表的方法
func (s *EndpointSet) ToDocs() (docs []interface{}) {
	for i := range s.Endpoints {
		docs = append(docs, s.Endpoints[i])
	}
	return
}

func NewRegistryResponse() *RegistryResponse {
	return &RegistryResponse{}
}
