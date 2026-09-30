package user

import (
	"github.com/radik-alt/auth_service/user/internal/repository"
	"github.com/radik-alt/auth_service/user/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo   repository.UserRepository
	bcryptCost int
}

func NewUserService(userRepo repository.UserRepository) service.UserService {
	return &userService{
		userRepo:   userRepo,
		bcryptCost: bcrypt.DefaultCost,
	}
}
