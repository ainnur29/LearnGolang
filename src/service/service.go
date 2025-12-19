package service

import (
	"learngolang/src/repository"
	"learngolang/src/service/user"
)

type Service struct {
	User user.UserServiceItf
}

func InitService(repository *repository.Repository) *Service {
	return &Service{
		User: user.InitUserService(
			repository.User,
		),
	}
}
