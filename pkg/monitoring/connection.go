package monitoring

import (
	"context"
	"errors"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const timeout = 10 * time.Second

var ErrSDKDisabled = errors.New("otel sdk disabled")
var ErrFailedToConnectToOtelAgent = errors.New("failed to connect to otel agent")

// Connect generates a new connection for opentelemetry.
func Connect(ctx context.Context) (*grpc.ClientConn, error) {
	if viper.GetBool("OTEL_SDK_DISABLED") {
		return nil, ErrSDKDisabled
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	otelAgentAddr := viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT")

	conn, err := grpc.DialContext(ctx, otelAgentAddr,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, errors.Join(ErrFailedToConnectToOtelAgent, err)
	}

	return conn, nil
}
