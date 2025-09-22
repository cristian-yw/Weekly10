package handlers

import (
	"net/http"
	"strconv"

	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	repo *repository.OrderRepository
}

func NewOrderHandler(repo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{repo: repo}
}

// @Summary     Get Movie Schedules with Filters
// @Description Get all schedules for a specific movie with optional filters:
// @Description - cinemaName (string, partial match)
// @Description - locationName (string, partial match)
// @Description - startTime (HH:MM format)
// @Description - date (YYYY-MM-DD format)
// @Tags        Orders
// @Produce     json
// @Param       movieId     path   int    true  "Movie ID"
// @Param       cinemaName  query  string false "Filter by cinema name (partial match)"
// @Param       locationName query string false "Filter by location name (partial match)"
// @Param       startTime   query  string false "Filter by start time (HH:MM)"
// @Param       date        query  string false "Filter by date (YYYY-MM-DD)"
// @Success     200 {array} models.Schedule
// @Failure     400 {object} map[string]string
// @Failure     500 {object} map[string]string
// @Security    BearerAuth
// @Router      /orders/{movieId}/schedules [get]
func (h *OrderHandler) GetSchedule(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movieId"})
		return
	}

	// Ambil query param filter (bisa kosong jika tidak diisi)
	cinemaName := c.Query("cinemaName")
	locationName := c.Query("locationName")
	startTime := c.Query("startTime")
	date := c.Query("date")

	schedules, err := h.repo.GetSchedule(
		c,
		movieID,
		cinemaName,
		locationName,
		startTime,
		date,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

// @Summary Get Available Seats
// @Description Get available seats for a specific schedule
// @Tags Orders
// @Produce json
// @Param scheduleId path int true "Schedule ID"
// @Success 200 {array} models.Seat
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /orders/seats/{scheduleId} [get]
func (h *OrderHandler) GetAvailableSeats(c *gin.Context) {
	scheduleID, _ := strconv.Atoi(c.Param("scheduleId"))
	seats, err := h.repo.GetAvailableSeats(c, scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seats)
}

// @Summary Get Movie Detail
// @Description Get detailed information of a specific movie
// @Tags Orders
// @Produce json
// @Param movieId path int true "Movie ID"
// @Success 200 {object} models.MovieDetail
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /orders/{movieId} [get]
func (h *OrderHandler) GetMovieDetail(c *gin.Context) {
	movieIDStr := c.Param("movieId")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid movie ID",
		})
		return
	}

	movie, err := h.repo.GetMovieDetail(c, movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// @Summary Create Order
// @Description Create a new order including seats selection
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body models.Order true "Order Request"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /orders/ [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req struct {
		UserID     int      `json:"user_id"`
		ScheduleID int      `json:"schedule_id"`
		TotalPrice int      `json:"total_price"`
		Seats      []string `json:"seats"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.repo.CreateOrder(c, req.UserID, req.ScheduleID, req.TotalPrice, req.Seats)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}
