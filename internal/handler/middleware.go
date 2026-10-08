package handler

import (
	"crypto/subtle"
	"log/slog"
	"net/http"
	"url-short/internal/config"

	"github.com/gin-gonic/gin"
)

func BasicAuth(cfg config.Auth, logger *slog.Logger) gin.HandlerFunc {
	return func(g *gin.Context) {
		user, pass, ok := g.Request.BasicAuth()
		if !ok {
			logger.Warn("нет авторизации", slog.String("ip", g.ClientIP()))
			g.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			g.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userOK := subtle.ConstantTimeCompare([]byte(user), []byte(cfg.User)) == 1
		passOK := subtle.ConstantTimeCompare([]byte(pass), []byte(cfg.Pass)) == 1

		if !userOK || !passOK {
			logger.Warn("неверные креды",
				slog.String("ip", g.ClientIP()),
				slog.String("user", user))
			g.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			g.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		g.Next()
	}
}
