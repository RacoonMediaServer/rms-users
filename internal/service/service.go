package service

import (
	"context"
	"fmt"
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"github.com/RacoonMediaServer/rms-users/internal/db"
	"github.com/RacoonMediaServer/rms-users/internal/model"
	"go-micro.dev/v4/logger"
)

type Service struct {
	db db.Users
}

func (s Service) CanAccess(ctx context.Context, request *rms_users.CanAccessRequest, response *rms_users.CanAccessResponse) error {
	u, err := s.db.FindUser(request.Token)
	if err != nil {
		logger.Errorf("attempt to find user failed: %s", err)
		return nil
	}
	if u == nil {
		logger.Warnf("user not found: %s", request.Token)
		return nil
	}

	switch request.Action {
	case rms_users.CanAccessRequest_Search:
		fallthrough
	case rms_users.CanAccessRequest_ConnectingToTheBot:
		response.Result = true
	case rms_users.CanAccessRequest_AccountManagement:
		response.Result = u.Admin
	}

	return nil
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
