package requests

type LoginRequest struct {
	Email    string `json:"email" validate:"required" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required" example:"11111111"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Name     string `json:"name" validate:"required" example:"John Doe"`
	Password string `json:"password" validate:"required,password" example:"11111111"`
}

type RefreshRequest struct {
	Token string `json:"token" validate:"required" example:"access_token"`
}
