package middleware

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type traceIdMiddlewareHandler struct{}

func (traceIdMiddlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(TraceIdKey)
	w.Write([]byte(traceId.(string)))
}

func TestSetTraceId(t *testing.T) {
	var err error
	var traceId string
	var updatedCtx context.Context
	var ctx context.Context

	var init = func() {
		ctx = context.Background()
	}

	// Test with invalid uuid
	init()
	traceId = "invalid trace"
	updatedCtx, err = SetTraceId(traceId, ctx)

	assert.NotNil(t, err)
	assert.Equal(t, ctx, updatedCtx)

	// Test with valid uuid
	init()
	traceId = uuid.NewString()
	updatedCtx, err = SetTraceId(traceId, ctx)

	assert.Nil(t, err)
	assert.Equal(t, traceId, updatedCtx.Value(TraceIdKey))
}

func TestGetTraceId(t *testing.T) {
	var err error
	var traceId string
	var ctx context.Context

	var init = func() {
		ctx = context.Background()
	}

	var initWithValue = func(value any) {
		init()
		ctx = context.WithValue(ctx, TraceIdKey, value)
	}

	// Test with undefined trace id
	init()
	traceId, err = GetTraceId(ctx)

	assert.Equal(t, UndefinedValue, err.Error())
	assert.Equal(t, "", traceId)

	// Test with unexpected value
	initWithValue(420)
	traceId, err = GetTraceId(ctx)

	assert.Equal(t, InvalidTraceId, err.Error())
	assert.Equal(t, "", traceId)

	// Test with valid uuit
	validTraceId := uuid.New().String()
	initWithValue(validTraceId)

	traceId, err = GetTraceId(ctx)

	assert.Nil(t, err)
	assert.Equal(t, validTraceId, traceId)
}

func TestNewChain(t *testing.T) {
	first := func(h http.Handler) http.Handler {
		return nil
	}

	second := func(h http.Handler) http.Handler {
		return nil
	}

	chain := NewChain(first, second)

	expected := reflect.ValueOf(first)
	actual := reflect.ValueOf(chain.items[0])
	assert.Equal(t, expected.Pointer(), actual.Pointer())

	expected = reflect.ValueOf(first)
	actual = reflect.ValueOf(chain.items[0])
	assert.Equal(t, expected.Pointer(), actual.Pointer())
}

func TestTraceMiddleware(t *testing.T) {
	server := httptest.NewServer(TraceMiddleware(traceIdMiddlewareHandler{}))

	res, err := http.Get(server.URL)
	assert.Nil(t, err)

	body, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	_, err = uuid.Parse(string(body))
	assert.Nil(t, err)
}
