package middleware

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func BuilChain(f http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) == 0 {
		return f
	}
	return m[0](BuilChain(f, m[1:cap(m)]...))
}


