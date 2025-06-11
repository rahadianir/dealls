package user

import (
	"context"

	"github.com/rahadianir/dealls/internal/models"
)

type UserRepositoryInterface interface {
	GetUserDetailsByUsername(ctx context.Context, username string) (models.User, error)
	GetUserRolesbyID(ctx context.Context, userID string) ([]string, error)
	GetAdminRole(ctx context.Context) (string, error)
	IsAdmin(ctx context.Context, userID string) (bool, error)
	GetUsersSalaryByIDs(ctx context.Context, userIDs []string) ([]models.UserSalary, error)
}

type UserLogicInterface interface {
	Login(ctx context.Context, username string, password string) (LoginResponse, error)
}
