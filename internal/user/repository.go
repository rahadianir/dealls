package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/huandu/go-sqlbuilder"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/models"
	"github.com/rahadianir/dealls/internal/pkg/xerror"
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
	sq := sqlbuilder.NewSelectBuilder()
	sq.Select(`id`, `name`, `username`, `password`, `salary`, `created_at`, `updated_at`, `deleted_at`, `created_by`, `updated_by`).
		From(`hr.users`).
		Where(
			sq.And(
				sq.Equal(`username`, username),
				sq.IsNull(`deleted_at`),
			),
		)
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var sqlUser SQLUser
	err := repo.deps.DB.QueryRowxContext(ctx, q, args...).StructScan(&sqlUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, xerror.ErrDataNotFound
		}

		return models.User{}, err
	}

	user := models.User{
		ID:        sqlUser.ID.String,
		Name:      sqlUser.Name.String,
		Username:  sqlUser.Username.String,
		Password:  sqlUser.Password.String,
		Salary:    sqlUser.Salary.Float64,
		CreatedAt: sqlUser.CreatedAt.Time,
		UpdatedAt: &sqlUser.DeletedAt.Time,
		DeletedAt: &sqlUser.DeletedAt.Time,
		CreatedBy: sqlUser.CreatedBy.String,
		UpdatedBy: sqlUser.UpdatedBy.String,
	}

	return user, nil
}

func (repo *UserRepository) GetUserRolesbyID(ctx context.Context, userID string) ([]string, error) {
	sq := sqlbuilder.NewSelectBuilder()
	sq.Select(`role_id`).
		From(`hr.user_role_map`).
		Where(
			sq.And(
				sq.Equal(`user_id`, userID),
				sq.IsNull(`deleted_at`),
			),
		)
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var result []string
	rows, err := repo.deps.DB.QueryxContext(ctx, q, args...)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	var temp string
	for rows.Next() {
		rows.Scan(&temp)
		result = append(result, temp)
	}

	if len(result) == 0 {
		return result, xerror.ErrDataNotFound
	}

	return result, nil
}

func (repo *UserRepository) GetAdminRole(ctx context.Context) (string, error) {
	sq := sqlbuilder.NewSelectBuilder()
	sq.Select(`id`).
		From(`hr.roles`).
		Where(
			sq.And(
				sq.Equal(`name`, `admin`),
				sq.IsNull(`deleted_at`),
			),
		)
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var result string
	err := repo.deps.DB.QueryRowxContext(ctx, q, args...).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, xerror.ErrDataNotFound
		}

		return result, err
	}

	return result, nil
}
