package service

import (
	"context"
	"errors"
	"fmt"
	"math"

	rms_bot_server "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-bot-server"
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/RacoonMediaServer/rms-users/internal/config"
	"github.com/RacoonMediaServer/rms-users/internal/model"
	jwt "github.com/golang-jwt/jwt/v5"
	"go-micro.dev/v4/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

var ErrUserNotFound = errors.New("user not found")

type Service struct {
	db Database
	f  servicemgr.ServiceFactory
}

// CheckPermissions implements rms_users.RmsUsersHandler.
func (s Service) CheckPermissions(ctx context.Context, req *rms_users.CheckPermissionsRequest, resp *rms_users.CheckPermissionsResponse) error {
	resp.Allowed = false

	claims := authClaims{}
	_, err := jwt.ParseWithClaims(req.Token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.Config().Security.Key), nil
	})
	if err != nil {
		unknownDeviceRequestsCounter.Inc()
		logger.Errorf("Auth failed [token='%s']: %s", req.Token, err)
		return nil
	}

	u, err := s.db.FindUser(claims.UserID)
	if err != nil {
		logger.Errorf("Attempt to find user failed: %s", err)
		return nil
	}
	if u == nil {
		unknownDeviceRequestsCounter.Inc()
		logger.Warnf("User not found: %s", claims.UserID)
		return nil
	}
	resp.UserId = claims.UserID

	for _, perm := range req.Perms {
		if !u.IsAllowed(perm) {
			return nil
		}
	}
	resp.Allowed = true
	deviceRequestsCounter.WithLabelValues(u.ID).Inc()
	return nil
}

func (s Service) RegisterUser(ctx context.Context, user *rms_users.User, response *rms_users.RegisterUserResponse) error {
	if user.TelegramUserID != nil {
		u, err := s.db.FindUserByTelegramID(*user.TelegramUserID)
		if err != nil {
			return err
		}
		if u != nil {
			token, err := s.GenerateAccessToken(u.ID)
			if err != nil {
				logger.Errorf("Generate access token failed: %s", err)
				return errors.New("error during add user")
			}
			response.Token = token
			response.UserId = u.ID
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

	accessToken, err := s.GenerateAccessToken(u.ID)
	if err != nil {
		logger.Errorf("Sign jwt token failed: %s", err)
		return errors.New("error during adding user")
	}

	response.Token = accessToken
	response.UserId = u.ID
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
		Id:             &u.ID,
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

func (s Service) CreateUser(user *model.User) (string, error) {
	if user.ID == "" {
		user.GenerateID()
	}
	token, err := s.GenerateAccessToken(user.ID)
	if err != nil {
		return "", err
	}
	if err = s.db.CreateUser(user); err != nil {
		return "", err
	}
	return token, err
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

func (s Service) IsAdminUser(token string) bool {
	resp := rms_users.CheckPermissionsResponse{}
	err := s.CheckPermissions(context.Background(), &rms_users.CheckPermissionsRequest{Token: token, Perms: []rms_users.Permissions{rms_users.Permissions_AccountManagement}}, &resp)
	if err != nil {
		return false
	}
	return resp.Allowed
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
	token, err := s.GenerateAccessToken(u.ID)
	if err != nil {
		logger.Errorf("Generate access token for admin failed: %s", err)
		return err
	}
	logger.Infof("Default admin key generated: %s", token)
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
				Id:             &u.ID,
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
