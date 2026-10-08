package handler

import (
	"log/slog"
	"url-short/internal/config"

	ginslog "github.com/gin-contrib/slog"
	"github.com/gin-gonic/gin"
)

// Router создаёт и настраивает gin.Engine со всеми маршрутами.
func (h *Handler) Router(authCfg config.Auth, logger *slog.Logger) *gin.Engine {
	r := gin.New()

	// свой логгер вместо дефолтного TextHandler из stderr
	r.Use(ginslog.SetLogger(ginslog.WithLogger(func(_ *gin.Context, _ *slog.Logger) *slog.Logger {
		return logger
	})))

	// паника → лог в slog + 500
	r.Use(gin.CustomRecovery(func(g *gin.Context, rec any) {
		logger.Error("panic recovered",
			slog.Any("panic", rec),
			slog.String("path", g.Request.URL.Path),
			slog.String("method", g.Request.Method))
		g.AbortWithStatusJSON(500, gin.H{"error": "internal error"})
	}))

	// Публичные маршруты — без auth
	r.GET("/:alias", h.Redirect)

	// Защищённые маршруты — с Basic Auth
	protected := r.Group("/")
	protected.Use(BasicAuth(authCfg, logger))
	{
		protected.GET("/urls/:alias", h.GetAlias)
		protected.POST("/urls", h.CreateAlias)
		protected.PATCH("/urls/:alias", h.UpdateAlias)
		protected.DELETE("/urls/:alias", h.DeleteAlias)
	}

	return r
}
