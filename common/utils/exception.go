package utils

import "github.com/infraboard/mcube/exception"

// IsAccessTokenExpiredError 判断是否是NotFoundError
func IsAccessTokenExpiredError(err error) bool {
	if err == nil {
		return false
	}

	e, ok := err.(exception.APIException)
	if !ok {
		return false
	}

	return e.ErrorCode() == exception.AccessTokenExpired && e.Namespace() == exception.GlobalNamespace.String()
}
