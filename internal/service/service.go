package service

import (
	"context"
	"errors"
	"fmt"
	rms_bot_server "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-bot-server"
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/RacoonMediaServer/rms-users/internal/model"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
)

var ErrUserNotFound = errors.New("user not found")

type Service struct {
	db Database
	f  servicemgr.ServiceFactory
}

func (s Service) RegisterUser(ctx context.Context, user *rms_users.User, response *rms_users.RegisterUserResponse) error {
	if user.TelegramUserID != nil {
		u, err := s.db.FindUserByTelegramID(*user.TelegramUserID)
		if err != nil {
			return err
		}
		if u != nil {
			for _, perm := range user.Perms {
				u.Grant(perm)
			}
			return s.db.UpdateUser(u)
		}
	}

	u := &model.User{
		Name:           &user.Name,
		TelegramUserId: user.TelegramUserID,
	}
	u.GenerateID()
	u.SetPermissions(user.Perms)
	response.Token = u.ID
	return s.db.CreateUser(u)
}

func (s Service) GetUserByTelegramId(ctx context.Context, request *rms_users.GetUserByTelegramIdRequest, user *rms_users.User) error {
	u, err := s.db.FindUserByTelegramID(request.TelegramUserId)
	if err != nil {
		return err
	}
	if u == nil {
		return nil
	}
	*user = rms_users.User{
		Token:          &u.ID,
		TelegramUserID: u.TelegramUserId,
		Perms:          u.GetPermissions(),
	}
	if u.Name != nil {
		user.Name = *u.Name
	}

	return nil
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

	if _, err = s.f.NewBotServer().DropSession(context.Background(), &rms_bot_server.DropSessionRequest{Token: ID}); err != nil {
		logger.Warnf("Notify bot failed: %s", err)
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
		unknownDeviceRequestsCounter.Inc()
		logger.Warnf("user not found: %s", request.Token)
		return nil
	}

	response.Perms = u.GetPermissions()
	deviceRequestsCounter.WithLabelValues(request.Token).Inc()

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

	return u.IsAllowed(rms_users.Permissions_AccountManagement)
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
		Info:        "Default admin user",
		Permissions: math.MaxInt32,
	}
	u.GenerateID()
	if err = s.db.CreateUser(&u); err != nil {
		return fmt.Errorf("store new admin user failed: %w", err)
	}
	logger.Infof("Default admin key generated: %s", u.ID)
	return nil
}

func (s Service) GetAdminUsers(ctx context.Context, empty *emptypb.Empty, response *rms_users.GetAdminUsersResponse) error {
	users, err := s.db.GetUsers()
	if err != nil {
		return err
	}
	for _, u := range users {
		if u.IsAllowed(rms_users.Permissions_AccountManagement) {
			result := &rms_users.User{
				Token:          &u.ID,
				TelegramUserID: u.TelegramUserId,
				Perms:          u.GetPermissions(),
			}
			response.Users = append(response.Users, result)
		}
	}
	return nil
}

func New(database Database, f servicemgr.ServiceFactory) Service {
	return Service{
		db: database,
		f:  f,
	}
}
