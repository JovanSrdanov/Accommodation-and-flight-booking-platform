package dto

type AccountRegistration struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" bindong:"required, email"`
}