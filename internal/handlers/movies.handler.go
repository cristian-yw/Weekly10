package handlers

import (
	"net/http"
	"os"
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

// @Summary Sync Popular Movies
// @Description Fetch popular movies from TMDB and store in database
// @Tags Movies
// @Produce json
// @Success 200 {object} models.SuccessMessage
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/sync/popular [post]
func (h *MovieHandler) SyncPopular(c *gin.Context) {
	apiKey := os.Getenv("API_KEY")
	if err := h.mr.SyncPopular(apiKey); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessMessage{Message: "Popular movies synced successfully"})
}

// @Summary Get Upcoming Movies
// @Description List upcoming movies with pagination
// @Tags Movies
// @Produce json
// @Param limit query int false "Number of movies per page" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.TMDBMovie
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/upcoming [get]
func (h *MovieHandler) GetUpcomingMovies(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	movies, err := h.mr.GetUpcomingMovies(ctx, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

// @Summary Get Popular Movies
// @Description List popular movies with pagination
// @Tags Movies
// @Produce json
// @Param limit query int false "Number of movies per page" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.TMDBMovie
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/popular [get]
func (h *MovieHandler) GetPopularMovies(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	movies, err := h.mr.GetPopularMovies(ctx, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, movies)
}

// @Summary Filter Movies
// @Description Filter movies by name or genre with pagination
// @Tags Movies
// @Produce json
// @Param name query string false "Movie name keyword"
// @Param genre_id query int false "Genre ID"
// @Param limit query int false "Number of movies per page" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.TMDBMovie
// @Failure 500 {object} models.ErrorResponse
// @Router /movies/filter [get]
func (h *MovieHandler) GetMoviesWithFilter(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	name := ctx.Query("name")
	genreID, _ := strconv.Atoi(ctx.DefaultQuery("genre_id", "0"))

	movies, err := h.mr.GetMoviesWithFilter(ctx, name, genreID, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, movies)
}
