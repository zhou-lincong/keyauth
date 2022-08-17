package utils_test

import (
	"testing"

	"github.com/zhou-lincong/keyauth/common/utils"
)

func TestToken(t *testing.T) {
	v := utils.MakeBearer(24)
	t.Log(v)
	// 	=== RUN   TestToken
	//     e:\goproject\keyauth\common\utils\token_test.go:11: FETFJcuKWV8iuG28RlgzmZCf
	// --- PASS: TestToken (0.00s)
	// PASS
	// ok  	github.com/zhou-lincong/keyauth/common/utils	0.671s

}
