package dto

type AccountRegistration struct {
	Username string `json:"username" binding:"required" bson:"username"`
	Password string `json:"password" binding:"required" bson:"password"`
	Email    string `json:"email" bindong:"required, email" bson:"email"`
}