package service

import (
	"github.com/Thalerngsak/restaurant/model"
	"github.com/Thalerngsak/restaurant/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (u userService) GetByID(userID uint) (*model.User, error) {
	return u.repo.GetByID(userID)
}

func (u userService) GetByUserName(userName string) (*model.User, error) {
	return u.repo.GetByUserName(userName)
}
