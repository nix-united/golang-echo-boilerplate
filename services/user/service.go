package user

import (
	"echo-demo-project/requests"

	"github.com/jinzhu/gorm"
)

type ServiceWrapper interface {
	Register(request *requests.RegisterRequest) error
}

type Service struct {
	Db *gorm.DB
}

func NewUserService(db *gorm.DB) *Service {
	return &Service{Db: db}
}
