package handler

import (
	"github.com/palanoyz/cinemahub/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	repo *repository.MovieRepository
}

func NewMovieHandler(repo *repository.MovieRepository) *MovieHandler {
	return &MovieHandler{repo: repo}
}

func (h *MovieHandler) List(c *gin.Context) {
	movies, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	c.JSON(http.StatusOK, movies)
}
