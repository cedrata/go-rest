package middleware

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

type args struct {
	context.Context
}

func TestSetTraceId(t *testing.T) {
	ctx := context.Background()

    traceId := uuid.NewString()
    ctx = SetTraceId(traceId, ctx)

    if res := ctx.Value(TraceIdKey); res != traceId {
        t.Fatalf("wrong traceId value\nprovided:%s\t returned:%s", traceId, res)
    }
}
