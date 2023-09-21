package middleware

import (
	"log"
	"net/http"
)

type ChainItem func(http.Handler) http.Handler

type Chain struct {
	items []ChainItem
}

func NewChain(items ...ChainItem) Chain {
	return Chain{append([]ChainItem(nil), items...)}
}

func (c Chain) Handle(h http.Handler) http.Handler {
	if h == nil {
	    	h = http.DefaultServeMux
	}

	for i := range c.items {
		h = c.items[len(c.items)-1-i](h)
	}

	return h
}

func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("start log")
		h.ServeHTTP(w, r)
		log.Println("end log")
	})
}

func HelloMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("start hello")
		h.ServeHTTP(w, r)
		log.Println("end log")
	})
}
