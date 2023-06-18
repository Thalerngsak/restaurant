package repository

import (
	"github.com/Thalerngsak/restaurant/model"
)

type UserRepository interface {
	GetByID(id uint) (*model.User, error)
	GetByUserName(userName string) (*model.User, error)
}
