package service

import (
	"context"
	"errors"
	"fmt"
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"github.com/RacoonMediaServer/rms-users/internal/db"
	"github.com/RacoonMediaServer/rms-users/internal/model"
	"go-micro.dev/v4/logger"
)

var ErrUserNotFound = errors.New("user not found")

type Service struct {
	db db.Users
}

func (s Service) GetUsers() ([]model.User, error) {
	return s.db.GetUsers()
}

func (s Service) CreateUser(user *model.User) error {
	if user.ID == "" {
		user.GenerateID()
	}
	return s.db.CreateUser(user)
}

func (s Service) DeleteUser(ID string) error {
	ok, err := s.db.DeleteUser(ID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrUserNotFound
	}
	return nil
}

func (s Service) GetPermissions(ctx context.Context, request *rms_users.GetPermissionsRequest, response *rms_users.GetPermissionsResponse) error {
	u, err := s.db.FindUser(request.Token)
	if err != nil {
		logger.Errorf("attempt to find user failed: %s", err)
		return nil
	}
	if u == nil {
		logger.Warnf("user not found: %s", request.Token)
		return nil
	}

	response.Perms = append(response.Perms, rms_users.Permissions_Search)
	response.Perms = append(response.Perms, rms_users.Permissions_ConnectingToTheBot)
	if u.Admin {
		response.Perms = append(response.Perms, rms_users.Permissions_AccountManagement)
	}

	return nil
}

func (s Service) IsAdminUser(ID string) bool {
	u, err := s.db.FindUser(ID)
	if err != nil {
		logger.Errorf("attempt to find user failed: %s", err)
		return false
	}
	if u == nil {
		logger.Warnf("user not found: %s", ID)
		return false
	}

	return u.Admin
}

func (s Service) CreateAdminIfNecessary() error {
	count, err := s.db.UsersCount()
	if err != nil {
		return fmt.Errorf("request count of users failed: %w", err)
	}
	if count != 0 {
		return nil
	}

	u := model.User{
		Info:  "Default admin user",
		Admin: true,
	}
	u.GenerateID()
	if err = s.db.CreateUser(&u); err != nil {
		return fmt.Errorf("store new admin user failed: %w", err)
	}
	logger.Infof("Default admin key generated: %s", u.ID)
	return nil
}

func New(database db.Users) Service {
	return Service{db: database}
}
