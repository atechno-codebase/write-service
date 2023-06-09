package logger

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func WithTracingFields(ctx context.Context) *logrus.Entry {

	span := trace.SpanFromContext(ctx)
	traceId := span.SpanContext().TraceID().String()
	spanId := span.SpanContext().SpanID().String()

	return logrus.
		WithField("trace_id", traceId).
		WithField("span_id", spanId)
}

func Info(ctx context.Context, args ...any) {

}
