package audit

import "time"

const (
	AppName = "audit"
)

type Service interface {
	RPCServer
}

func NewOperateLog(who, resource, action string) *OperateLog {
	return &OperateLog{
		Username: who,
		When:     time.Now().UnixMilli(),
		Resource: resource,
		Action:   action,
	}
}
