package handler

import (
	"url-short/internal/config"

	ginslog "github.com/gin-contrib/slog"
	"github.com/gin-gonic/gin"
)

// Router создаёт и настраивает gin.Engine со всеми маршрутами.
func (h *Handler) Router(authCfg config.Auth) *gin.Engine {
	r := gin.New() // пустой движок, без встроенных middleware

	// Глобальные middleware — применяются ко всем запросам
	r.Use(ginslog.SetLogger()) // логирование через slog
	r.Use(gin.Recovery())      // восстановление после паник → 500

	// Публичные маршруты — без auth
	r.GET("/:alias", h.Redirect)

	// Защищённые маршруты — с Basic Auth
	protected := r.Group("/")
	protected.Use(BasicAuth(authCfg))
	{
		protected.GET("/urls/:alias", h.GetAlias)
		protected.POST("/urls", h.CreateAlias)
		protected.PATCH("/urls/:alias", h.UpdateAlias)
		protected.DELETE("/urls/:alias", h.DeleteAlias)
	}

	return r
}
