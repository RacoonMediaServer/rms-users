package server

import (
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"github.com/RacoonMediaServer/rms-users/internal/server/restapi/operations/registration"
	"github.com/go-openapi/runtime/middleware"
	"go-micro.dev/v4/logger"
)

func (s *Server) signUp(params registration.SignUpParams) middleware.Responder {
	if !s.Registration.Enabled {
		logger.Warn("Self registration is disabled")
		return registration.NewSignUpForbidden()
	}

	found := false
	for _, allowedDomain := range s.Registration.Domains {
		if params.Domain == allowedDomain {
			found = true
			break
		}
	}

	if !found {
		logger.Warnf("Attempt to sign up to denied domain: %s", params.Domain)
		return registration.NewSignUpForbidden()
	}

	user := rms_users.User{
		Perms: []rms_users.Permissions{
			rms_users.Permissions_ConnectingToTheBot,
			rms_users.Permissions_Search,
		},
		Domain: &params.Domain,
	}

	resp := rms_users.RegisterUserResponse{}

	if err := s.Users.RegisterUser(params.HTTPRequest.Context(), &user, &resp); err != nil {
		logger.Errorf("Register user failed: %s", err)
		return registration.NewSignUpInternalServerError()
	}

	logger.Infof("User registered by self-registration: domain = %s, id = %s", params.Domain, resp.UserId)
	return registration.NewSignUpOK().WithPayload(&registration.SignUpOKBody{Token: &resp.Token})
}
