package service

import (
	"github.com/Thalerngsak/restaurant/model"
)

type UserService interface {
	GetByID(userID uint) (*model.User, error)
	GetByUserName(userName string) (*model.User, error)
}
