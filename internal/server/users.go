package server

import (
	"errors"
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
		payload.Results = append(payload.Results, &users.GetUsersOKBodyResultsItems0{
			ID:      user.ID,
			Info:    user.Info,
			IsAdmin: user.Admin,
		})
	}
	return users.NewGetUsersOK().WithPayload(&payload)
}

func (s *Server) createUser(params users.CreateUserParams, key *models.Principal) middleware.Responder {
	isAdmin := false
	if params.User.IsAdmin != nil {
		isAdmin = *params.User.IsAdmin
	}
	u := model.User{
		Info:  *params.User.Info,
		Admin: isAdmin,
	}
	err := s.Users.CreateUser(&u)
	if err != nil {
		logger.Errorf("Create user failed: %s", err)
		return users.NewCreateUserInternalServerError()
	}

	return users.NewCreateUserOK().WithPayload(&users.CreateUserOKBody{ID: &u.ID})
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
