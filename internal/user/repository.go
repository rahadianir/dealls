package user

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/huandu/go-sqlbuilder"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/models"
	"github.com/rahadianir/dealls/internal/pkg/dbhelper"
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

func (repo *UserRepository) IsAdmin(ctx context.Context, userID string) (bool, error) {
	// get user roles
	roleIDs, err := repo.GetUserRolesbyID(ctx, userID)
	if err != nil {
		if errors.Is(err, xerror.ErrDataNotFound) {
			repo.deps.Logger.WarnContext(ctx, "user has no role", slog.Any("error", err))
			return false, nil
		}
		repo.deps.Logger.ErrorContext(ctx, "failed to get user roles by id", slog.Any("error", err))
		return false, err
	}

	// get admin role
	adminRoleID, err := repo.GetAdminRole(ctx)
	if err != nil {
		repo.deps.Logger.ErrorContext(ctx, "failed to fetch admin role ID", slog.Any("error", err))
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

func (repo *UserRepository) GetUsersSalaryByIDs(ctx context.Context, userIDs []string) ([]models.UserSalary, error) {
	sq := sqlbuilder.NewSelectBuilder()
	sq.Select(`id`, `salary`).
		From(`hr.users`).
		Where(
			sq.And(
				sq.In(`id::text`, sqlbuilder.List(userIDs)),
				sq.IsNull(`deleted_at`),
			),
		)
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)

	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil {
		return []models.UserSalary{}, err
	}

	var temp SQLUserSalary
	var result []models.UserSalary
	for rows.Next() {
		err := rows.StructScan(&temp)
		if err != nil {
			repo.deps.Logger.WarnContext(ctx, "failed to scan user salary", slog.Any("error", err))
		}

		result = append(result, models.UserSalary{
			UserID: temp.ID.String,
			Salary: temp.Salary.Float64,
		})
	}

	return result, nil
}
