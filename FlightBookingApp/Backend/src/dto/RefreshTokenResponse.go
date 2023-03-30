package dto

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken" binding:"required"`
}
