package user

import (
	"context"

	"github.com/rahadianir/dealls/internal/config"
	"golang.org/x/crypto/bcrypt"
)

type UserLogic struct {
	deps     *config.CommonDependencies
	userRepo UserRepository
}

func NewUserLogic(deps *config.CommonDependencies, userRepo UserRepository) *UserLogic {
	return &UserLogic{
		deps:     deps,
		userRepo: userRepo,
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

	}

	
	return "", nil
}

func (logic *UserLogic) IsAdmin(ctx context.Context, userID string) (bool, error) {
	return true, nil
}
