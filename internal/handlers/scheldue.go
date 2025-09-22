package handlers

import (
	"net/http"

	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	repo *repository.ScheduleRepository
}

func NewScheduleHandler(repo *repository.ScheduleRepository) *ScheduleHandler {
	return &ScheduleHandler{repo: repo}
}

func (h *ScheduleHandler) GetGenresHandler(c *gin.Context) {
	genres, err := h.repo.GetGenres(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, genres)
}

func (h *ScheduleHandler) GetCinemasHandler(c *gin.Context) {
	cinemas, err := h.repo.GetCinemas(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cinemas)
}

func (h *ScheduleHandler) GetLocationsHandler(c *gin.Context) {
	locations, err := h.repo.GetLocations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locations)
}
