package middleware

import (
	"context"
	"log"
	"net/http"
	"reflect"

	"github.com/google/uuid"
)

type ContextKey string

const ContextReqId ContextKey = "ReqId"

// add function to get id from context
// add function to set id for provided contex

type ChainItem func(http.Handler) http.Handler

type Chain struct {
	items []ChainItem
}

func NewChain(items ...ChainItem) Chain {
	return Chain{append([]ChainItem(nil), items...)}
}

// in function to create new chain automatically add log handler as first one.

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
        ctx := context.WithValue(r.Context(), ContextReqId, uuid.New().String())
        log.Println("generated ", reflect.TypeOf(ctx.Value(ContextReqId)))
		h.ServeHTTP(w, r.WithContext(ctx))
		log.Println("end log")
	})
}

func HelloMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("start hello")
		log.Println(r.Context().Value(ContextReqId))
		h.ServeHTTP(w, r)
		log.Println("end log")
	})
}
