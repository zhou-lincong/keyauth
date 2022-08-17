package utils

import (
	"math/rand"
	"strings"
	"time"
)

// MakeBearer https://tools.ietf.org/html/rfc6750#section-2.1
// 要求是b64编码的token =
// 只能编码这些字符：1*( ALPHA / DIGIT /"-" / "." / "_" / "~" / "+" / "/" ) *"="
// 只起到填充作用
func MakeBearer(lenth int) string {
	charlist := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	t := make([]string, lenth)
	rand.Seed(time.Now().UnixNano() + int64(lenth) + rand.Int63n(10000))
	for i := 0; i < lenth; i++ {
		rn := rand.Intn(len(charlist))
		w := charlist[rn : rn+1]
		t = append(t, w)
	}

	token := strings.Join(t, "")
	return token
}
