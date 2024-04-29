package appcontext

import (
	"context"

	"github.com/google/uuid"
)

type ContextKey string

const (
	TraceIDKey  ContextKey = "trace_id"
	UserIDKey   ContextKey = "user_id"
	TenantIDKey ContextKey = "tenant_id"
)

type Context interface {
	context.Context
	UserID() (string, bool)
	SetUserID(userID string) Context
	TenantID() (string, bool)
	SetTenantID(tenantID string) Context
	TraceID() string
	SetTraceID(traceID string) Context
}

// AppContext carries the context of the current execution.
type appContext struct {
	// original context
	context.Context
}

// UserID returns the user id
func (sc *appContext) UserID() (string, bool) {
	userID := sc.Context.Value(UserIDKey)

	if id, ok := userID.(string); ok {
		return id, ok
	}

	return "", false
}

// SetUserID sets the user id
func (sc *appContext) SetUserID(userID string) Context {
	sc.Context = context.WithValue(sc.Context, UserIDKey, userID)

	return sc
}

// TenantID returns the tenant id
func (sc *appContext) TenantID() (string, bool) {
	tenantIDKey := sc.Context.Value(TenantIDKey)

	if id, ok := tenantIDKey.(string); ok {
		return id, true
	}

	return "", false
}

// SetTenantID sets the user id
func (sc *appContext) SetTenantID(tenantID string) Context {
	sc.Context = context.WithValue(sc.Context, TenantIDKey, tenantID)

	return sc
}

// SetTraceID sets the trace id
func (sc *appContext) SetTraceID(traceID string) Context {
	sc.Context = context.WithValue(sc.Context, TraceIDKey, traceID)

	return sc
}

// TraceID returns the trace identifier for the current flow
func (sc *appContext) TraceID() string {
	return sc.Context.Value(TraceIDKey).(string)
}

// FromContext returns a new AppContext from a context.Context
func FromContext(ctx context.Context) Context {
	appCtx := NewAppContext(ctx)

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		appCtx.SetTraceID(traceID)
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		appCtx.SetUserID(userID)
	}

	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok {
		appCtx.SetTenantID(tenantID)
	}

	return appCtx
}

// NewContext returns a new AppContext
func NewAppContext(ctx context.Context) Context {
	ctx = context.WithValue(ctx, TraceIDKey, uuid.NewString())
	return &appContext{Context: ctx}
}
