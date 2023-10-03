package middleware

import (
	"context"
	"log"
	"net/http"
	"reflect"

	"github.com/google/uuid"
)

type ContextKey string
type ContextValue string

const TraceIdKey ContextKey = "ReqId"
const UndefinedValue ContextValue = "undefined"
const InvalidTraceId ContextValue = "invalid trace id value in context"

// Set a trace id in the given context.
func SetTraceId(value string, ctx context.Context) context.Context {
	return context.WithValue(ctx, TraceIdKey, value)
}

// Access to the trace id given key with type safety.
// If the variable is not defined or an error comes up retriving it 
// an according value is returned.
func GetTraceId(ctx context.Context) ContextValue {
	value := ctx.Value(TraceIdKey)

	if value == nil {
		return UndefinedValue
	}
	traceId, ok := value.(ContextValue)

	if !ok {
		return InvalidTraceId
	}

	return traceId
}

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
		ctx := context.WithValue(r.Context(), TraceIdKey, uuid.New().String())
		log.Println("generated ", reflect.TypeOf(ctx.Value(TraceIdKey)))
		h.ServeHTTP(w, r.WithContext(ctx))
		log.Println("end log")
	})
}

func HelloMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("start hello")
		log.Println(r.Context().Value(TraceIdKey))
		h.ServeHTTP(w, r)
		log.Println("end log")
	})
}
