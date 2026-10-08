package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Redirect(g *gin.Context) {
	alias := g.Param("alias")

	url, err := h.svc.GetUrlByAlias(g.Request.Context(), alias)
	if err != nil {
		writeErr(g, err)
		return
	}

	g.Redirect(http.StatusFound, url.OriginalURL)
}
