package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/cristian-yw/Weekly10/internal/models"
	"github.com/cristian-yw/Weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
// @Security BearerAuth
// @Router /user/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("userID")
	profile, err := h.repo.GetByID(c.Request.Context(), userID)
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
// @Security BearerAuth
// @Router /user/history [get]
func (h *UserHandler) GetHistory(c *gin.Context) {
	userID := c.GetInt("userID")
	history, err := h.repo.GetHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Edit profile sekaligus upload avatar (opsional). Hanya field yang dikirim yang akan di-update.
// @Tags         Users
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        first_name formData string false "First name"
// @Param        last_name  formData string false "Last name"
// @Param        phone      formData string false "Phone number"
// @Param        avatar     formData file   false "Avatar image"
// @Success      200 {object} models.User
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /user/profile [patch]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt("userID")

	var firstName, lastName, phone *string
	if v := c.PostForm("first_name"); v != "" {
		firstName = &v
	}
	if v := c.PostForm("last_name"); v != "" {
		lastName = &v
	}
	if v := c.PostForm("phone"); v != "" {
		phone = &v
	}

	// optional avatar file
	var avatarURL *string
	if file, err := c.FormFile("avatar"); err == nil {
		filename := fmt.Sprintf("avatar_%d_%d%s", userID, time.Now().Unix(), filepath.Ext(file.Filename))
		path := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(file, path); err == nil {
			u := "/uploads/" + filename
			avatarURL = &u
		}
	}

	if err := h.repo.UpdateProfile(
		c.Request.Context(),
		userID,
		firstName,
		lastName,
		phone,
		avatarURL,
	); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.repo.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// ChangePassword godoc
// @Summary      Change user password
// @Description  Ganti password akun. User harus mengirim password lama dan password baru.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.ChangePasswordRequest  true  "Current & new password"
// @Success      200   {object}  map[string]string             "message: password updated successfully"
// @Failure      400   {object}  map[string]string             "invalid request body"
// @Failure      401   {object}  map[string]string             "current password is incorrect"
// @Failure      500   {object}  map[string]string             "internal server error"
// @Router       /user/password [patch]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt("userID")

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// 1. Ambil hash password lama langsung dari tabel users
	storedHash, err := h.repo.GetPasswordHash(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get password"})
		return
	}

	// 2. Validasi password lama
	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.CurrentPassword)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
		return
	}

	// 3. Hash password baru
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// 4. Simpan password baru
	if err := h.repo.UpdatePassword(c.Request.Context(), userID, string(hashed)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}
