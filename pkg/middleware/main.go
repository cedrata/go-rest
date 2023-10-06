package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"reflect"

	"github.com/google/uuid"
)

type ChainItem func(http.Handler) http.Handler

type Chain struct {
	items []ChainItem
}

const (
	TraceIdKey     string = "ReqId"
	UndefinedValue string = "undefined value"
	InvalidTraceId string = "invalid trace id value in context"
)

// Set a trace id in the given context.
// If the provided traceId is not a valid uuid the same context
// is returned with some error.
func SetTraceId(value string, ctx context.Context) (context.Context, error) {
	_, err := uuid.Parse(value)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, TraceIdKey, value), nil
}

// Access to the trace id given key with type safety.
// If the variable is not defined or an error comes up retriving it
// an according value is returned.
func GetTraceId(ctx context.Context) (string, error) {
	value := ctx.Value(TraceIdKey)

	if value == nil {
		return "", errors.New(UndefinedValue)
	}

	traceId, ok := value.(string)

	if !ok {
		return "", errors.New(InvalidTraceId)
	}

	return traceId, nil
}

// Create a new middleware chain with the given middleware functions.
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

func TraceMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), TraceIdKey, uuid.New().String())
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func LogMiddleware(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("start log")
// 		ctx := context.WithValue(r.Context(), TraceIdKey, uuid.New().String())
// 		log.Println("generated ", reflect.TypeOf(ctx.Value(TraceIdKey)))
// 		h.ServeHTTP(w, r.WithContext(ctx))
// 		log.Println("end log")
// 	})
// }

// func HelloMiddleware(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("start hello")
// 		log.Println(r.Context().Value(TraceIdKey))
// h.ServeHTTP(w, r)
// 		log.Println("end log")
// 	})
// }
