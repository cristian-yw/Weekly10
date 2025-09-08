package models

type MovieRequest struct {
	Title   string `json:"title"`
	Runtime int    `json:"runtime"`
}
type SuccessMessage struct {
	Message string `json:"message" example:"Movie updated successfully"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"something went wrong"`
}
