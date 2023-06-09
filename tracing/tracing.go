package tracing

import "go.opentelemetry.io/otel/trace"

var tracing trace.Tracer

func init() {

}

func Tracer() trace.Tracer {
	return trace
}
