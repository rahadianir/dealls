package user

import (
	"context"

	"github.com/rahadianir/dealls/internal/config"
)

type UserLogic struct {
	deps *config.CommonDependencies
}

func NewUserLogic(deps *config.CommonDependencies) *UserLogic {
	return &UserLogic{
		deps: deps,
	}
}

func (logic *UserLogic) Login(ctx context.Context, username string, password string) (string, error) {

	return "", nil
}

func (logic *UserLogic) IsAdmin(ctx context.Context, userID string) (bool, error) {
	return true, nil
}
