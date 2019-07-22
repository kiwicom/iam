package monitoring

import (
	"context"
	"net"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// TracerOptions are required options to create an instance of a tracer
type TracerOptions struct {
	ServiceName string
	Environment string
	Host        string
	Port        string
}

// Tracer is the Tracer instance that contains all the methods
type Tracer struct {
}

// CreateNewTracingService creates a new Instance of Tracer
func CreateNewTracingService(options TracerOptions) (*Tracer, error) {
	addr := net.JoinHostPort(options.Host, options.Port)
	tracer.Start(
		tracer.WithAgentAddr(addr),
		tracer.WithServiceName(options.ServiceName),
		tracer.WithGlobalTag("env", options.Environment),
	)
	return &Tracer{}, nil
}

// StartSpan creates a new tracing span and returns it
func (*Tracer) StartSpan(operationName, resourceName, resourceType string) tracer.Span {
	span := tracer.StartSpan(operationName, tracer.ResourceName(resourceName), tracer.SpanType(resourceType))
	return span
}

// StartSpanWithContext creates a new span attached to current context
func (*Tracer) StartSpanWithContext(
	oldContext context.Context,
	operationName,
	resourceName,
	resourceType string) (tracer.Span, context.Context) {
	span, newContext := tracer.StartSpanFromContext(
		oldContext,
		operationName,
		tracer.ResourceName(resourceName),
		tracer.SpanType(resourceType),
	)
	return span, newContext
}

// FinishSpan finishes tracking given span.
// Should be used in conjunction with defer to trigger at the end of every function that should be tracked
func (*Tracer) FinishSpan(span tracer.Span) {
	span.Finish()
}

// Stop stops the instance of Tracer
func (*Tracer) Stop() {
	tracer.Stop()
}
