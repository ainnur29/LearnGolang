package service

import (
	"golang-bulang-bolang/src/repository"
	"golang-bulang-bolang/src/service/user"
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
