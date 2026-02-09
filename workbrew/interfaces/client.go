package interfaces

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// HTTPClient interface that services will use
// This breaks import cycles by providing a contract without implementation
type HTTPClient interface {
	Get(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) error
	Post(ctx context.Context, path string, body any, headers map[string]string, result any) error
	PostWithQuery(ctx context.Context, path string, queryParams map[string]string, body any, headers map[string]string, result any) error
	Put(ctx context.Context, path string, body any, headers map[string]string, result any) error
	Patch(ctx context.Context, path string, body any, headers map[string]string, result any) error
	Delete(ctx context.Context, path string, queryParams map[string]string, headers map[string]string, result any) error
	DeleteWithBody(ctx context.Context, path string, body any, headers map[string]string, result any) error
	GetCSV(ctx context.Context, path string, queryParams map[string]string, headers map[string]string) ([]byte, error)
	GetLogger() *zap.Logger
	QueryBuilder() ServiceQueryBuilder
}

// ServiceQueryBuilder defines the query builder contract for services
type ServiceQueryBuilder interface {
	AddString(key, value string) QueryBuilder
	AddInt(key string, value int) QueryBuilder
	AddInt64(key string, value int64) QueryBuilder
	AddBool(key string, value bool) QueryBuilder
	AddTime(key string, value time.Time) QueryBuilder
	AddStringSlice(key string, values []string) QueryBuilder
	AddIntSlice(key string, values []int) QueryBuilder
	AddCustom(key, value string) QueryBuilder
	AddIfNotEmpty(key, value string) QueryBuilder
	AddIfTrue(condition bool, key, value string) QueryBuilder
	Merge(other map[string]string) QueryBuilder
	Remove(key string) QueryBuilder
	Has(key string) bool
	Get(key string) string
	Build() map[string]string
	BuildString() string
	Clear() QueryBuilder
	Count() int
	IsEmpty() bool
}

// QueryBuilder interface for method chaining
type QueryBuilder interface {
	ServiceQueryBuilder
}
