package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/pkg/xerror"
	"github.com/rahadianir/dealls/internal/pkg/xjwt"
	"golang.org/x/crypto/bcrypt"
)

type UserLogic struct {
	deps      *config.CommonDependencies
	userRepo  UserRepository
	jwtHelper xjwt.JWTHelper
}

func NewUserLogic(deps *config.CommonDependencies, userRepo UserRepository, jwtHelper xjwt.JWTHelper) *UserLogic {
	return &UserLogic{
		deps:      deps,
		userRepo:  userRepo,
		jwtHelper: jwtHelper,
	}
}

func (logic *UserLogic) Login(ctx context.Context, username string, password string) (string, error) {
	userDetails, err := logic.userRepo.GetUserDetailsByUsername(ctx, username)
	if err != nil {
		// do something
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(password))
	if err != nil {

		return "", err
	}

	token, err := logic.jwtHelper.GenerateJWT(logic.deps.Config.App.Name, userDetails.ID, logic.deps.Config.App.ExpiryTime, logic.deps.Config.App.SecretKey)
	if err != nil {
		return token, err
	}

	return token, nil
}

func (logic *UserLogic) IsAdmin(ctx context.Context, userID string) (bool, error) {
	// get user roles
	roleIDs, err := logic.userRepo.GetUserRolesbyID(ctx, userID)
	if err != nil {
		if errors.Is(err, xerror.ErrDataNotFound) {
			logic.deps.Logger.WarnContext(ctx, "user has no role", slog.Any("error", err))
			return false, nil
		}

		return false, err
	}

	// get admin role
	adminRoleID, err := logic.userRepo.GetAdminRole(ctx)
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to fetch admin role ID", slog.Any("error", err))
		return false, err
	}

	// match user roles with admin role
	for _, id := range roleIDs {
		if id == adminRoleID {
			return true, nil
		}
	}

	return false, nil
}
