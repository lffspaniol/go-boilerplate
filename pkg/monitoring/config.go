package monitoring

import (
	"context"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
)

// Setup start the resource for identfy the app.
func Setup() (*resource.Resource, error) {
	res, err := resource.New(
		context.Background(),
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(viper.GetString("SERVICE_NAME")),
			semconv.DeploymentEnvironmentKey.String("production"),
		),
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
