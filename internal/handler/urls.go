package handler

import (
	"errors"
	"net/http"
	"url-short/internal/domain"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	OriginalURL string `json:"url" binding:"required"`
	Alias       string `json:"alias" binding:"required"`
}

func wtiterErr(g *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		g.JSON(http.StatusNotFound, gin.H{"error": err})
	case errors.Is(err, domain.ErrConflict):
		g.JSON(http.StatusConflict, gin.H{"error": err})
	case errors.Is(err, domain.ErrValidation):
		g.JSON(http.StatusBadRequest, gin.H{"error": err})
	default:
		// неизвестная ошибка (БД упала) — наружу не светим
		g.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}

func (h *Handler) CreateAlias(g *gin.Context) {
	var req CreateRequest

	err := g.ShouldBindJSON(&req)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "неправильный запрос"})
		return
	}

	err = h.svc.SaveUrl(g.Request.Context(), req.OriginalURL, req.Alias)
	if err != nil {
		wtiterErr(g, err)
		return
	}

	g.JSON(201, gin.H{"status": "created"})
}

func (h *Handler) DeleteAlias(g *gin.Context) {

	par := g.Param("alias")

	err := h.svc.DeleteUrlByAlias(g.Request.Context(), par)
	if err != nil {
		wtiterErr(g, err)
		return
	}

	g.JSON(200, gin.H{"status": "deleted"})
}

func (h *Handler) UpdateAlias(g *gin.Context) {
	var req CreateRequest

	err := g.ShouldBindJSON(&req)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "неправильный запрос"})
		return
	}

	err = h.svc.UpdateUrlByAlias(g.Request.Context(), req.Alias, req.OriginalURL)
	if err != nil {
		wtiterErr(g, err)
		return
	}

	g.JSON(200, gin.H{"status": "update"})
}

func (h *Handler) GetAlias(g *gin.Context) {

	par := g.Param("alias")

	d, err := h.svc.GetUrlByAlias(g.Request.Context(), par)
	if err != nil {
		wtiterErr(g, err)
		return
	}

	g.JSON(200, gin.H{"id": d.ID,
		"original_url": d.OriginalURL,
		"alias":        d.Alias,
		"created_at":   d.CreatedAt,
		"updated_at":   d.UpdateAt})
}
