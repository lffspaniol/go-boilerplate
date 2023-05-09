package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type Shutdown func(ctx context.Context) error

// NewTracer creates a new open telemetry tracer and returns a shutdown function.
func NewTracer(
	ctx context.Context,
	conn *grpc.ClientConn,
	res *resource.Resource,
	logger *zap.Logger,
) (Shutdown, error) {
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	batch := sdktrace.NewBatchSpanProcessor(traceExporter)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(batch),
	)

	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func(ctx context.Context) error {
		if shutErr := tracerProvider.Shutdown(ctx); shutErr != nil {
			logger.Error("failed to shutdown tracer provider", zap.Error(shutErr))
			return err
		}
		return nil
	}, nil
}
