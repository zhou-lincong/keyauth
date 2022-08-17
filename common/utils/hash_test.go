package utils_test

import (
	"testing"

	"github.com/zhou-lincong/keyauth/common/utils"
)

// a9993e364706816aba3e25717850c26c9cd0d89d
// 但这种方式每次hash的结果都是一样
func TestHash(t *testing.T) {
	v := utils.Hash("abc")
	t.Log(v)
}

//第一次：$2a$14$/joRw3xC/A5Bc0HrcmBP1OB1cXgo2Pj/SwBbdWVlhQDS2ZhncV1bO
//第二次：$2a$14$XutsTRs4fwici5XPyxsOUOi3.y0Jn6TOKM4cFhuHau68a33MHVcxm
func TestPasswordHash(t *testing.T) {
	v := utils.HashPassword("abc")
	t.Log(v)
	ok := utils.CheckPasswordHash("abc", "$2a$14$/joRw3xC/A5Bc0HrcmBP1OB1cXgo2Pj/SwBbdWVlhQDS2ZhncV1bO")
	t.Log(ok)
	ok = utils.CheckPasswordHash("abc", "$2a$14$XutsTRs4fwici5XPyxsOUOi3.y0Jn6TOKM4cFhuHau68a33MHVcxm")
	t.Log(ok)

	// 	=== RUN   TestPasswordHash
	//     e:\goproject\keyauth\common\utils\hash_test.go:20: $2a$14$6zl8HgWLlZ5m1lXouiEUZu5drZAfo4sAHebf26XFkZJ5wG6WQGpfa
	//     e:\goproject\keyauth\common\utils\hash_test.go:22: true
	//     e:\goproject\keyauth\common\utils\hash_test.go:24: true
	// --- PASS: TestPasswordHash (7.00s)
	// PASS
	// ok  	github.com/zhou-lincong/keyauth/common/utils	7.553s
}
