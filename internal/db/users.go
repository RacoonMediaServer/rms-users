package db

import (
	"errors"
	"github.com/RacoonMediaServer/rms-users/internal/model"
	"gorm.io/gorm"
)

// Users is an interface for user management
type Users interface {
	// UsersCount returns total users
	UsersCount() (uint64, error)

	// CreateUser appends user to DB
	CreateUser(user *model.User) error

	// FindUser finds user by ID (nil means not found)
	FindUser(ID string) (*model.User, error)

	// GetUsers returns all users records
	GetUsers() ([]model.User, error)

	// DeleteUser deletes user and returns affected rows
	DeleteUser(ID string) (bool, error)
}

func (d database) UsersCount() (uint64, error) {
	var count int64
	result := d.conn.Model(&model.User{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return uint64(count), nil
}

func (d database) CreateUser(user *model.User) error {
	return d.conn.Create(user).Error
}

func (d database) FindUser(ID string) (*model.User, error) {
	var u model.User
	if err := d.conn.First(&u, "id = ?", ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (d database) GetUsers() ([]model.User, error) {
	result := make([]model.User, 0)
	if err := d.conn.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (d database) DeleteUser(ID string) (bool, error) {
	tx := d.conn.Model(&model.User{}).Unscoped().Delete(&model.User{ID: ID})
	if tx.Error != nil {
		return false, tx.Error
	}

	return tx.RowsAffected != 0, nil
}
