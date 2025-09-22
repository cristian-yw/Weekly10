package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	mr *repository.MovieRepository
}

func NewMovieHandler(mr *repository.MovieRepository) *MovieHandler {
	return &MovieHandler{mr: mr}
}

// @Summary Get Upcoming Movies
// @Tags Movies
// @Produce json
// @Param limit query int false "Number of movies per page" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.TMDBMovie
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/upcoming [get]
func (h *MovieHandler) GetUpcomingMovies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	movies, err := h.mr.GetUpcomingMovies(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

// @Summary Get Popular Movies
// @Tags Movies
// @Produce json
// @Param limit query int false "Number of movies per page" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.TMDBMovie
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/popular [get]
func (h *MovieHandler) GetPopularMovies(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	movies, err := h.mr.GetPopularMovies(c, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

// @Summary Filter Movies
// @Tags Movies
// @Produce json
// @Param name query string false "Movie name keyword"
// @Param genre_id query int false "Genre ID"
// @Param limit query int false "Number of movies per page" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.TMDBMovie
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/filter [get]
func (h *MovieHandler) GetMoviesWithFilter(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	name := c.Query("name")
	genreID, _ := strconv.Atoi(c.DefaultQuery("genre_id", "0"))

	// 1. Ambil movies dengan filter & pagination
	movies, err := h.mr.GetMoviesWithFilter(c, name, genreID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 2. Hitung total record untuk pagination
	totalRecords, err := h.mr.CountMoviesWithFilter(c, name, genreID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	c.JSON(http.StatusOK, gin.H{
		"results":     movies,
		"total_pages": totalPages,
		"total_items": totalRecords,
	})
}

// @Summary Get all movies
// @Tags Movies
// @Produce json
// @Success 200 {array} models.TMDBMovie
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/all [get]
func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	movies, err := h.mr.GetAllMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to fetch movies"})
		return
	}
	c.JSON(http.StatusOK, movies)
}
