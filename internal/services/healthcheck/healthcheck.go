package healthcheck

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"
)

const pkgName = "healthcheck"
const OK = "OK"

var ErrDepencenciesFailed = errors.New("depencencies failed")

type Depencencie interface {
	Healthcheck(context.Context) error
}

type Alive struct {
	depencencies []Depencencie
}

func (alive *Alive) Readiness(ctx context.Context) error {
	ctx, span := otel.Tracer(pkgName).Start(ctx, "Readiness")
	defer span.End()

	g, ctx := errgroup.WithContext(ctx)

	for _, depencency := range alive.depencencies {
		depencency := depencency
		g.Go(func() error {
			return depencency.Healthcheck(ctx)
		})
	}

	if err := g.Wait(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.Join(ErrDepencenciesFailed, err)
	}

	return nil
}

func New(depencencies ...Depencencie) Alive {
	return Alive{
		depencencies: depencencies,
	}
}
