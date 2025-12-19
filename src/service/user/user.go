package user

import (
	"context"

	"learngolang/src/domain"
	"learngolang/src/repository/user"
)

type UserServiceItf interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error)
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	ListUsers(ctx context.Context, filter domain.UserFilter) ([]*domain.User, int64, error)
	UpdateUser(ctx context.Context, id int64, req domain.UpdateUserRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type userService struct {
	userRepository user.UserRepositoryItf
}

func InitUserService(userRepository user.UserRepositoryItf) UserServiceItf {
	return &userService{
		userRepository: userRepository,
	}
}
