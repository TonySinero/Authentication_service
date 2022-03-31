package model

type AuthUser struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"password"`
}

type RestorePassword struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string
}
