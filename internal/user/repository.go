package user

import (
	"context"

	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/models"
)

type UserRepository struct {
	deps *config.CommonDependencies
}

func NewUserRepository(deps *config.CommonDependencies) *UserRepository {
	return &UserRepository{
		deps: deps,
	}
}

func (repo *UserRepository) GetUserDetailsByUsername(ctx context.Context, username string) (models.User, error) {
	
	return models.User{}, nil
}
