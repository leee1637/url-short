package handler

import (
	"errors"
	"net/http"
	"url-short/internal/domain"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	OriginalURL string `json:"url" binding:"required"`
	Alias       string `json:"alias"`
}

type UpdateRequest struct {
	OriginalURL string `json:"url" binding:"required"`
}

func writeErr(g *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		g.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrConflict):
		g.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrValidation):
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		g.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
	}
}

func (h *Handler) CreateAlias(g *gin.Context) {
	var req CreateRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err := h.svc.SaveUrl(g.Request.Context(), req.OriginalURL, req.Alias); err != nil {
		writeErr(g, err)
		return
	}

	g.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (h *Handler) DeleteAlias(g *gin.Context) {
	alias := g.Param("alias")

	if err := h.svc.DeleteUrlByAlias(g.Request.Context(), alias); err != nil {
		writeErr(g, err)
		return
	}

	g.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *Handler) UpdateAlias(g *gin.Context) {
	var req UpdateRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	alias := g.Param("alias")

	if err := h.svc.UpdateUrlByAlias(g.Request.Context(), alias, req.OriginalURL); err != nil {
		writeErr(g, err)
		return
	}

	g.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *Handler) GetAlias(g *gin.Context) {
	alias := g.Param("alias")

	d, err := h.svc.GetUrlByAlias(g.Request.Context(), alias)
	if err != nil {
		writeErr(g, err)
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"id":           d.ID,
		"original_url": d.OriginalURL,
		"alias":        d.Alias,
		"created_at":   d.CreatedAt,
		"updated_at":   d.UpdateAt,
	})
}
