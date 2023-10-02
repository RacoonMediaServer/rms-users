package db

import (
	"errors"
	"github.com/RacoonMediaServer/rms-users/internal/model"
	"gorm.io/gorm"
)

func (d Database) UsersCount() (uint64, error) {
	var count int64
	result := d.conn.Model(&model.User{}).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return uint64(count), nil
}

func (d Database) CreateUser(user *model.User) error {
	return d.conn.Create(user).Error
}

func (d Database) FindUser(ID string) (*model.User, error) {
	var u model.User
	if err := d.conn.First(&u, "id = ?", ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (d Database) FindUserByTelegramID(ID int32) (*model.User, error) {
	var u model.User
	if err := d.conn.First(&u, "telegram_user_id = ?", ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (d Database) GetUsers() ([]model.User, error) {
	result := make([]model.User, 0)
	if err := d.conn.Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (d Database) DeleteUser(ID string) (bool, error) {
	tx := d.conn.Model(&model.User{}).Unscoped().Delete(&model.User{ID: ID})
	if tx.Error != nil {
		return false, tx.Error
	}

	return tx.RowsAffected != 0, nil
}

func (d Database) UpdateUser(u *model.User) error {
	return d.conn.Save(u).Error
}
