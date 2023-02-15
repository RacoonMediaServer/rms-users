package service

import (
	"context"
	rms_users "github.com/RacoonMediaServer/rms-packages/pkg/service/rms-users"
	"github.com/RacoonMediaServer/rms-users/internal/db"
)

type Service struct {
	db db.Users
}

func (s Service) CanAccess(ctx context.Context, request *rms_users.CanAccessRequest, response *rms_users.CanAccessResponse) error {
	//TODO implement me
	panic("implement me")
}

func New(database db.Users) Service {
	return Service{db: database}
}
