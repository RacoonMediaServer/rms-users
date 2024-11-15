package server

import (
	"errors"
	"math"

	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"github.com/RacoonMediaServer/rms-users/internal/model"
	"github.com/RacoonMediaServer/rms-users/internal/server/models"
	"github.com/RacoonMediaServer/rms-users/internal/server/restapi/operations/users"
	"github.com/RacoonMediaServer/rms-users/internal/service"
	"github.com/go-openapi/runtime/middleware"
	"go-micro.dev/v4/logger"
)

func (s *Server) getUsers(params users.GetUsersParams, key *models.Principal) middleware.Responder {
	registeredUsers, err := s.Users.GetUsers()
	if err != nil {
		logger.Errorf("Get users failed: %s", err)
		return users.NewGetUsersInternalServerError()
	}

	payload := users.GetUsersOKBody{}
	for _, user := range registeredUsers {
		u := &models.User{
			ID:     user.ID,
			Info:   user.Info,
			Domain: user.Domain,
		}
		if user.TelegramUserId != nil {
			u.TelegramUser = int64(*user.TelegramUserId)
		}
		if user.Name != nil {
			u.Name = *user.Name
		}
		u.Role = new(string)
		if user.IsAllowed(rms_users.Permissions_AccountManagement) {
			*u.Role = "admin"
		} else if user.IsAllowed(rms_users.Permissions_ConnectingToTheBot) {
			*u.Role = "user"
		} else if user.IsAllowed(rms_users.Permissions_ListeningMusic) {
			*u.Role = "listener"
		}
		payload.Results = append(payload.Results, u)
	}
	return users.NewGetUsersOK().WithPayload(&payload)
}

func (s *Server) createUser(params users.CreateUserParams, key *models.Principal) middleware.Responder {
	u := model.User{
		Name:   &params.User.Name,
		Info:   params.User.Info,
		Domain: params.User.Domain,
	}

	switch *params.User.Role {
	case "admin":
		u.Permissions = math.MaxInt32
	case "user":
		u.SetPermissions([]rms_users.Permissions{
			rms_users.Permissions_Search,
			rms_users.Permissions_ConnectingToTheBot,
			rms_users.Permissions_SendNotifications,
			rms_users.Permissions_ListeningMusic,
		})
	case "listener":
		u.Grant(rms_users.Permissions_Search)
		u.Grant(rms_users.Permissions_ListeningMusic)
	}

	token, err := s.Users.CreateUser(&u)
	if err != nil {
		logger.Errorf("Create user failed: %s", err)
		return users.NewCreateUserInternalServerError()
	}

	return users.NewCreateUserOK().WithPayload(&users.CreateUserOKBody{ID: &u.ID, Token: &token})
}

func (s *Server) deleteUser(params users.DeleteUserParams, key *models.Principal) middleware.Responder {
	if err := s.Users.DeleteUser(params.ID); err != nil {
		logger.Errorf("Delete user failed: %s", err)
		if errors.Is(err, service.ErrUserNotFound) {
			return users.NewDeleteUserNotFound()
		}
		return users.NewDeleteUserInternalServerError()
	}

	return users.NewDeleteUserOK()
}
