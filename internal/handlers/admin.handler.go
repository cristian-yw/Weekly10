package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

// ===================== CRUD =====================

// @Summary Create a new movie (with schedules)
// @Description Admin can add a new movie along with its schedules (cinema, location, time, date).
// @Tags Admin
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param title formData string true "Movie Title"
// @Param overview formData string true "Overview"
// @Param release_date formData string true "Release Date (YYYY-MM-DD)"
// @Param runtime formData int true "Runtime in minutes"
// @Param tmdb_id formData int false "TMDB ID"
// @Param popularity formData number false "Popularity"
// @Param vote_average formData number false "Vote Average"
// @Param vote_count formData int false "Vote Count"
// @Param genres formData string false "Comma separated genres (e.g. Action,Drama)"
// @Param schedules formData string false "JSON array of schedules. Example: [{\"cinema_id\":1,\"location_id\":1,\"time_id\":2,\"date\":\"2025-09-21\"}]"
// @Param poster formData file false "Poster image file"
// @Param backdrop formData file false "Backdrop image file"
// @Success 201 {object} map[string]interface{} "movie_id returned"
// @Failure 400 {object} models.ErrorResponse "Invalid input"
// @Failure 401 {object} models.ErrorResponse "Unauthorized (Missing or invalid token)"
// @Failure 500 {object} models.ErrorResponse "Server error"
// @Router /admin/movies [post]
func (h *AdminHandler) CreateMovie(c *gin.Context) {
	var req models.NewMovieRequest

	// Title & Overview
	req.Title = c.PostForm("title")
	req.Overview = c.PostForm("overview")
	req.ReleaseDate = c.PostForm("release_date")

	// Runtime
	if runtimeStr := c.PostForm("runtime"); runtimeStr != "" {
		if val, err := strconv.Atoi(runtimeStr); err == nil {
			req.Runtime = val
		}
	}

	// TMDB ID (nullable → pointer)
	if tmdbStr := c.PostForm("tmdb_id"); tmdbStr != "" {
		if val, err := strconv.Atoi(tmdbStr); err == nil {
			req.TMDBID = &val
		}
	} else {
		req.TMDBID = nil
	}

	// Popularity
	if popStr := c.PostForm("popularity"); popStr != "" {
		if val, err := strconv.ParseFloat(popStr, 64); err == nil {
			req.Popularity = val
		}
	}

	// Vote average
	if vaStr := c.PostForm("vote_average"); vaStr != "" {
		if val, err := strconv.ParseFloat(vaStr, 64); err == nil {
			req.VoteAverage = val
		}
	}

	// Vote count
	if vcStr := c.PostForm("vote_count"); vcStr != "" {
		if val, err := strconv.Atoi(vcStr); err == nil {
			req.VoteCount = val
		}
	}

	// Parse genres: ekspektasi "Action,Drama"
	if gnames := c.PostForm("genres"); gnames != "" {
		for _, s := range strings.Split(gnames, ",") {
			name := strings.TrimSpace(s)
			if name != "" {
				req.Genres = append(req.Genres, name)
			}
		}
	}

	// Parse schedules: JSON string
	if schStr := c.PostForm("schedules"); schStr != "" {
		if err := json.Unmarshal([]byte(schStr), &req.Schedules); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid schedules format"})
			return
		}
	}

	// Parse director: JSON object {"tmdb_id":123, "name":"Christopher Nolan"}
	if dirStr := c.PostForm("director"); dirStr != "" {
		var director models.PersonRequest
		if err := json.Unmarshal([]byte(dirStr), &director); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid director format"})
			return
		}
		req.Director = &director
	}

	// Parse casts: JSON array [{"tmdb_id":1,"name":"Leonardo DiCaprio"}]
	if castStr := c.PostForm("casts"); castStr != "" {
		var casts []models.PersonRequest
		if err := json.Unmarshal([]byte(castStr), &casts); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid casts format"})
			return
		}
		req.Casts = casts
	}

	// Upload poster
	if file, err := c.FormFile("poster"); err == nil {
		if err := os.MkdirAll("/", os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		dst := filepath.Join(fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename))
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		req.PosterPath = dst
	}

	// Upload backdrop
	if file, err := c.FormFile("backdrop"); err == nil {
		if err := os.MkdirAll("uploads/backdrop", os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		dst := filepath.Join("/", fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename))
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
			return
		}
		req.BackdropPath = dst
	}

	// Simpan ke DB
	movieID, err := h.repo.CreateMovie(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "movie created successfully",
		"movie_id": movieID,
	})
}

// @Summary Get movie by ID
// @Tags Admin
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.TMDBMovie
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /admin/movies/{id} [get]
func (h *AdminHandler) GetMovieByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid movie id"})
		return
	}

	movie, err := h.repo.GetMovieByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// @Summary Patch update movie
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param movie body object true "Partial movie update (only send fields you want to update)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /admin/movies/{id} [patch]
func (h *AdminHandler) PatchMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid movie id"})
		return
	}

	// Bind input ke map agar fleksibel
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Hapus field `id` biar tidak bisa diubah manual
	delete(input, "id")

	if len(input) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "no fields provided for update"})
		return
	}

	// Panggil repo untuk update dinamis
	err = h.repo.PatchMovie(c.Request.Context(), id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to patch movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie patched successfully"})
}

// @Summary Delete movie
// @Tags Admin
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /admin/movies/{id} [delete]
func (h *AdminHandler) DeleteMovie(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "invalid movie id"})
		return
	}

	err = h.repo.DeleteMovie(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "failed to delete movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie deleted successfully"})
}

// ===================== TMDB SYNC & LIST =====================

// @Summary Sync Popular Movies
// @Description Fetch popular movies from TMDB and store in database
// @Tags Admin
// @Produce json
// @Success 200 {object} models.SuccessMessage
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /admin/sync/popular [post]
func (h *AdminHandler) SyncPopular(c *gin.Context) {
	apiKey := os.Getenv("API_KEY")
	if err := h.repo.SyncPopular(apiKey); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessMessage{Message: "Popular movies synced successfully"})
}
