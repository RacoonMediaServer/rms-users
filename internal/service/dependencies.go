package service

import "github.com/RacoonMediaServer/rms-users/internal/model"

type Database interface {
	UsersCount() (uint64, error)
	CreateUser(user *model.User) error
	FindUser(ID string) (*model.User, error)
	FindUserByTelegramID(id int32) (*model.User, error)
	GetUsers() ([]model.User, error)
	DeleteUser(ID string) (bool, error)
	UpdateUser(user *model.User) error
}
