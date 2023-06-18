package repository

import (
	"github.com/Thalerngsak/restaurant/model"
)

type userStore struct {
}

func NewUserDB() UserRepository {
	return userStore{}
}

func (u userStore) GetByID(id uint) (*model.User, error) {
	user := model.User{ID: 1, Username: "thalerngsak"}
	return &user, nil
}

func (u userStore) GetByUserName(userName string) (*model.User, error) {
	user := model.User{ID: 1, Username: "thalerngsak"}
	return &user, nil
}
