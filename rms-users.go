package main

import (
	"fmt"
	"net/http"
	"sync"

	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"github.com/RacoonMediaServer/rms-packages/pkg/service/servicemgr"
	"github.com/RacoonMediaServer/rms-users/internal/config"
	"github.com/RacoonMediaServer/rms-users/internal/db"
	"github.com/RacoonMediaServer/rms-users/internal/server"
	userService "github.com/RacoonMediaServer/rms-users/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"

	// Plugins
	_ "github.com/go-micro/plugins/v4/registry/etcd"
)

var Version = "v0.0.0"

const serviceName = "rms-users"

func main() {
	logger.Infof("%s %s", serviceName, Version)
	defer logger.Info("DONE.")

	useDebug := false

	service := micro.NewService(
		micro.Name(serviceName),
		micro.Version(Version),
		micro.Flags(
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"debug"},
				Usage:       "debug log level",
				Value:       false,
				Destination: &useDebug,
			},
		),
	)

	service.Init(
		micro.Action(func(context *cli.Context) error {
			configFile := fmt.Sprintf("/etc/rms/%s.json", serviceName)
			if context.IsSet("config") {
				configFile = context.String("config")
			}
			return config.Load(configFile)
		}),
	)

	if useDebug {
		_ = logger.Init(logger.WithLevel(logger.DebugLevel))
	}

	database, err := db.Connect(config.Config().Database)
	if err != nil {
		logger.Fatalf("Connect to database failed: %s", err)
	}

	handler := userService.New(database, servicemgr.NewServiceFactory(service), config.Config().Security)

	// создаем пользователя-админа по умолчанию
	if err = handler.CreateAdminIfNecessary(); err != nil {
		logger.Fatalf("Create admin user failed: %s", err)
	}

	// регистрируем хендлеры
	if err := rms_users.RegisterRmsUsersHandler(service.Server(), handler); err != nil {
		logger.Fatalf("Register service failed: %s", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		srv := server.Server{Users: handler}
		if err := srv.ListenAndServer(config.Config().Http.Host, config.Config().Http.Port); err != nil {
			logger.Fatalf("Cannot start web server: %+s", err)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Config().Monitor.Host, config.Config().Monitor.Port), nil); err != nil {
			logger.Fatalf("Cannot bind monitoring endpoint: %s", err)
		}
	}()

	if err := service.Run(); err != nil {
		logger.Fatalf("Run service failed: %s", err)
	}

	wg.Wait()
}
