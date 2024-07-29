package main

import (
	"context"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// newTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func newTracerProvider(serviceName, serviceID, url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(version),
			semconv.ServiceInstanceIDKey.String(serviceID),
		)),
	)
	return tp, nil
}

func myTraceHandler() server.HandlerWrapper {
	// return a handler wrapper
	return func(h server.HandlerFunc) server.HandlerFunc {
		// return a function that returns a function
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			key := "X-Kardinal-Trace-Id"

			keyFromMD, found := metadata.Get(ctx, key)
			if found {
				logger.Infof("[LEO-DEBUG] key: %s", keyFromMD)
			} else {
				logger.Infof("[LEO-DEBUG] header key not found")
			}

			return nil
		}
	}
}
