package controler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

func (ctrl *Controlers) HandleEcho(c *gin.Context) {
	ctx, span := otel.Tracer(pkgName).Start(c.Request.Context(), "HandleEcho")
	defer span.End()
	response, err := ctrl.echo.Echo(ctx, "message")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		ctrl.log.Error("Echo failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, response)
}
