package entity

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
type SuccessResponse struct {
	Message string `json:"message" example:"success message"`
}
