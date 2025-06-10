package user

import (
	"context"
	"errors"
	"fmt"
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

func (logic *UserLogic) Login(ctx context.Context, username string, password string) (LoginResponse, error) {
	userDetails, err := logic.userRepo.GetUserDetailsByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, xerror.ErrDataNotFound) {
			logic.deps.Logger.WarnContext(ctx, "username not found", slog.Any("error", err))
			return LoginResponse{}, err
		}

		return LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(password))
	if err != nil {
		logic.deps.Logger.WarnContext(ctx, "invalid password", slog.Any("error", err))
		return LoginResponse{}, xerror.ClientError{Err: fmt.Errorf("invalid password")}
	}

	token, err := logic.jwtHelper.GenerateJWT(logic.deps.Config.App.Name, userDetails.ID, logic.deps.Config.App.ExpiryTime, logic.deps.Config.App.JWTSecretKey)
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to generate JWT", slog.Any("error", err))
		return LoginResponse{}, err
	}

	result := LoginResponse{
		Token: token,
	}

	if logic.deps.Config.App.IsDebug {
		result.UserID = userDetails.ID
	}

	return result, nil
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
