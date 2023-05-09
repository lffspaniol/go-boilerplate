package controler

import (
	"boilerplate/internal/services/echo"
	"boilerplate/internal/services/healthcheck"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

const pkgName = "controler"

type Controlers struct {
	alive healthcheck.Alive
	echo  echo.Service
	log   *zap.Logger
}

func (ctrl *Controlers) HandleHeathCheck(c *gin.Context) {
	c.JSON(http.StatusOK, healthcheck.OK)
}

func (ctrl *Controlers) HandleReadiness(c *gin.Context) {
	ctx, span := otel.Tracer(pkgName).Start(c.Request.Context(), "HandleReadiness")
	defer span.End()
	if err := ctrl.alive.Readiness(ctx); err != nil {
		ctrl.log.Error("Readiness failed", zap.Error(err))
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, healthcheck.OK)
}

func New(log *zap.Logger, alive healthcheck.Alive, echo echo.Service) *Controlers {
	return &Controlers{
		alive: alive,
		echo:  echo,
		log:   log,
	}
}
