package server

import (
	"fmt"
	"github.com/RacoonMediaServer/rms-users/internal/server/restapi"
	"github.com/RacoonMediaServer/rms-users/internal/server/restapi/operations"
	"github.com/RacoonMediaServer/rms-users/internal/service"
	"github.com/go-openapi/loads"
	"go-micro.dev/v4/logger"
)

type Server struct {
	srv *restapi.Server

	Users service.Service
}

func (s *Server) ListenAndServer(host string, port int) error {
	if s.srv == nil {
		swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
		if err != nil {
			return err
		}

		// создаем хендлеры API по умолчанию
		api := operations.NewServerAPI(swaggerSpec)
		s.configureAPI(api)

		// устанавливаем свой логгер
		api.Logger = func(content string, i ...interface{}) {
			logger.Infof(content, i...)
		}

		// создаем и настраиваем сервер
		s.srv = restapi.NewServer(api)
	}

	s.srv.Host = host
	s.srv.Port = port

	if err := s.srv.Listen(); err != nil {
		return fmt.Errorf("cannot start server: %w", err)
	}

	return s.srv.Serve()
}

func (s *Server) Shutdown() error {
	if s.srv != nil {
		return s.srv.Shutdown()
	}

	return nil
}
