package token

import (
	"fmt"
	"time"

	"github.com/infraboard/mcube/exception"
	"github.com/zhou-lincong/keyauth/apps/user"
	"github.com/zhou-lincong/keyauth/common/utils"
)

const (
	AppName = "token"
)

func NewIssueTokenRequest() *IssueTokenRequest {
	return &IssueTokenRequest{
		UserDomain: user.DefaultDomain,
	}
}

func (req *IssueTokenRequest) Validate() error {
	switch req.GranteType {
	case GranteType_PASSWORD:
		if req.Username == "" || req.Password == "" {
			return fmt.Errorf("password grant required username and password")
		}
	}

	return nil
}

func NewToken(req *IssueTokenRequest, expiredDuration time.Duration) *Token {
	now := time.Now()
	// Token 10
	expired := now.Add(expiredDuration)
	refresh := now.Add(expiredDuration * 5)

	return &Token{
		AccessToken:           utils.MakeBearer(24),
		IssueAt:               now.UnixMilli(),
		Data:                  req,
		AccessTokenExpiredAt:  expired.UnixMilli(),
		RefreshToken:          utils.MakeBearer(32),
		RefreshTokenExpiredAt: refresh.UnixMilli(),
	}
}

func NewDefaultToken() *Token {
	return &Token{
		Data: &IssueTokenRequest{},
		Meta: map[string]string{},
	}
}

// 判断Token过期没有
func (t *Token) Validate() error {
	// 是一个时间戳,
	//  now现在的时间   expire过期时间
	if time.Now().UnixMilli() > t.AccessTokenExpiredAt {
		return exception.NewAccessTokenExpired("access token expired")
	}

	return nil
}

// 判断refresh Token过期没有
func (t *Token) IsRefreshTokenExpired() bool {
	// 是一个时间戳,
	//  now   expire
	if time.Now().UnixMilli() > t.RefreshTokenExpiredAt {
		return true
	}

	return false
}

// 续约Token, 延长一个生命周期
func (t *Token) Extend(expiredDuration time.Duration) {
	now := time.Now()
	// Token 10
	expired := now.Add(expiredDuration)
	refresh := now.Add(expiredDuration * 5)

	t.AccessTokenExpiredAt = expired.UnixMilli()
	t.RefreshTokenExpiredAt = refresh.UnixMilli()
}

func NewDescribeTokenRequest(at string) *DescribeTokenRequest {
	return &DescribeTokenRequest{
		AccessToken: at,
	}
}

func NewValidateTokenRequest(at string) *ValidateTokenRequest {
	return &ValidateTokenRequest{
		AccessToken: at,
	}
}

func NewRevolkTokenRequest() *RevolkTokenRequest {
	return &RevolkTokenRequest{}
}
