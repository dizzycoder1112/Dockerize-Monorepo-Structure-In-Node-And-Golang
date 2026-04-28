package postgres

import (
	"context"
	"fmt"
	"time"

	logger "dizzycoder1112/Dockerize-Monorepo-Structure-In-Node-And-Golang/logger"
	"github.com/jackc/pgx/v5"
)

type contextKey string

const (
	queryStartKey contextKey = "db_query_start"
	querySQLKey   contextKey = "db_query_sql"
)

const slowQueryThresholdMs = 200.0

// QueryLogger implements pgx.QueryTracer. In dev/staging it logs every query
// with full SQL at Debug level; in production it logs only queries slower than
// slowQueryThresholdMs at Warn and omits parameter values.
type QueryLogger struct {
	log logger.Logger
	env string
}

func NewQueryLogger(log logger.Logger, env string) *QueryLogger {
	return &QueryLogger{log: log, env: env}
}

func (ql *QueryLogger) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	ctx = context.WithValue(ctx, queryStartKey, time.Now())
	ctx = context.WithValue(ctx, querySQLKey, data.SQL)
	return ctx
}

func (ql *QueryLogger) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	start, ok := ctx.Value(queryStartKey).(time.Time)
	if !ok {
		return
	}
	durationMs := float64(time.Since(start).Nanoseconds()) / 1e6
	sql, _ := ctx.Value(querySQLKey).(string)

	if ql.env == "production" {
		if durationMs < slowQueryThresholdMs {
			return
		}
		fields := []any{
			"duration_ms", fmt.Sprintf("%.2f", durationMs),
			"slow", true,
		}
		if data.Err != nil {
			fields = append(fields, "error", data.Err)
		}
		ql.log.Warn("db_query", fields...)
		return
	}

	fields := []any{
		"duration_ms", fmt.Sprintf("%.2f", durationMs),
		"sql", sql,
	}
	if data.Err != nil {
		fields = append(fields, "error", data.Err)
	}
	ql.log.Debug("db_query", fields...)
}
