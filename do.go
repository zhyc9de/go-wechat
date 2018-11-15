package weUtil

import "time"

type Token struct {
	Token   string
	Expires int64
}

func (token Token) IsNotExpire() bool {
	return token.Expires > time.Now().Unix()
}
