package payroll

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/rahadianir/dealls/internal/config"
)

type PayrollRepository struct {
	deps *config.CommonDependencies
}

func NewPayrollRepository(deps *config.CommonDependencies) *PayrollRepository {
	return &PayrollRepository{
		deps: deps,
	}
}

func (repo *PayrollRepository) SetPayrollPeriod(ctx context.Context, start time.Time, end time.Time) error {
	ins := sqlbuilder.NewInsertBuilder()
	insertQ, insertArgs := ins.InsertInto(`hr.attendance_period`).
		Cols(`id`, `start_date`, `end_date`, `active`, `created_at`).
		Values(uuid.NewString(), start, end, true).BuildWithFlavor(sqlbuilder.PostgreSQL)
	update := sqlbuilder.NewUpdateBuilder()
	updateQ, updateArgs := update.Update(`hr.attendance_period`).Set(update.Assign(`active`, false)).BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx, err := repo.deps.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, updateQ, updateArgs...)
	if err != nil {
		repo.deps.Logger.ErrorContext(ctx, "failed to set all payroll period inactive", slog.Any("error", err))
		return err
	}

	_, err = tx.ExecContext(ctx, insertQ, insertArgs...)
	if err != nil {
		repo.deps.Logger.ErrorContext(ctx, "failed to insert new payroll period", slog.Any("error", err))
		return err
	}

	err = tx.Commit()
	if err != nil {
		repo.deps.Logger.ErrorContext(ctx, "failed to commit new payroll period", slog.Any("error", err))
		return err
	}

	return nil
}
