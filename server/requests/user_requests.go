package requests

type LoginRequest struct {
	Name string  `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}


type RegisterRequest struct {
	Name string  `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}