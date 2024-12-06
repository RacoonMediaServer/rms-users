package server

import (
	"net/http"

	"github.com/RacoonMediaServer/rms-users/internal/server/models"
	"github.com/RacoonMediaServer/rms-users/internal/server/restapi/operations"
	"github.com/RacoonMediaServer/rms-users/internal/server/restapi/operations/registration"
	"github.com/RacoonMediaServer/rms-users/internal/server/restapi/operations/users"
	"github.com/go-openapi/errors"
	"go-micro.dev/v4/logger"
)

func (s *Server) configureAPI(api *operations.ServerAPI) {
	api.UsersGetUsersHandler = users.GetUsersHandlerFunc(s.getUsers)
	api.UsersCreateUserHandler = users.CreateUserHandlerFunc(s.createUser)
	api.UsersDeleteUserHandler = users.DeleteUserHandlerFunc(s.deleteUser)
	api.RegistrationSignUpHandler = registration.SignUpHandlerFunc(s.signUp)

	api.KeyAuth = func(token string) (*models.Principal, error) {
		if !s.Users.IsAdminUser(token) {
			logger.Warnf("Attempt of accessing to he restricted area: %s", token)
			return nil, errors.New(http.StatusForbidden, "Forbidden")
		}
		return &models.Principal{Token: token}, nil
	}
}
