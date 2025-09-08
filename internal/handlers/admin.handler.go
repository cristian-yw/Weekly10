package handlers

import (
	"net/http"
	"strconv"

	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	repo *repository.AdminRepository
}

func NewAdminHandler(repo *repository.AdminRepository) *AdminHandler {
	return &AdminHandler{repo: repo}
}

// @Summary Get All Movies
// @Description Get list of all movies (Admin only)
// @Tags Admin
// @Produce json
// @Success 200 {array} models.Movie
// @Failure 500 {object} models.ErrorResponse
// @Security Bearer
// @Router /admin/movies [get]
func (h *AdminHandler) GetAllMovies(c *gin.Context) {
	movies, err := h.repo.GetAllMovies(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

// @Summary Update Movie
// @Description Update movie by ID (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param request body models.MovieRequest true "Movie info"
// @Success 200 {object} models.SuccessMessage
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security Bearer
// @Router /admin/movies/{id} [put]
func (h *AdminHandler) UpdateMovie(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req models.MovieRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.repo.UpdateMovie(c, id, req.Title, req.Runtime); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessMessage{Message: "Movie updated successfully"})
}

// @Summary Delete Movie
// @Description Delete movie by ID (Admin only)
// @Tags Admin
// @Param id path int true "Movie ID"
// @Success 200 {object} models.SuccessMessage
// @Failure 500 {object} models.ErrorResponse
// @Security Bearer
// @Router /admin/movies/{id} [delete]
func (h *AdminHandler) DeleteMovie(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.repo.DeleteMovie(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessMessage{Message: "Movie deleted successfully"})
}
