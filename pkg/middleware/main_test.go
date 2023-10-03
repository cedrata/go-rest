package middleware

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

type args struct {
	context.Context
}

func TestSetTraceId(t *testing.T) {
    ctx := context.Background()

    SetTraceId(uuid.NewString(), ctx)

    fmt.Println(ctx.Value(TraceIdKey))
}
