package handlers

import (
	"net/http"

	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// @Summary Get User Profile
// @Description Get logged-in user's profile information
// @Tags Users
// @Produce json
// @Success 200 {object} models.UserProfileResponse
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /auth/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("userID")
	profile, err := h.repo.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

// @Summary Get Order History
// @Description Get logged-in user's order history
// @Tags Users
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /auth/history [get]
func (h *UserHandler) GetHistory(c *gin.Context) {
	userID := c.GetInt("userID")
	history, err := h.repo.GetHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

// @Summary Edit Profile
// @Description Update logged-in user's profile
// @Tags Users
// @Accept json
// @Produce json
// @Param request body models.UserProfileUpdate true "Profile update request"
// @Success 200 {object} models.UserProfileResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security Bearer
// @Router /auth/profile [put]
func (h *UserHandler) EditProfile(c *gin.Context) {
	userID := c.GetInt("userID")
	var req models.UserProfileUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	profile, err := h.repo.EditProfile(c.Request.Context(), userID, req.FirstName, req.LastName, req.Phone, req.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}
