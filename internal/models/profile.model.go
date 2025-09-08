package models

import "time"

type UserProfile struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
}

type UserProfileResponse struct {
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
	Role      string `json:"role"`
}

type UserProfileUpdate struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
}
type OrderHistory struct {
	OrderID    int       `json:"order_id"`
	MovieTitle string    `json:"movie_title"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	Date       time.Time `json:"date"`
}

type EditProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
