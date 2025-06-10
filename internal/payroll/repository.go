package payroll

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/pkg/dbhelper"
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
	insertQ, insertArgs := ins.InsertInto(`hr.attendance_periods`).
		Cols(`id`, `start_date`, `end_date`, `active`, `created_at`).
		Values(uuid.NewString(), start, end, true, `now()`).BuildWithFlavor(sqlbuilder.PostgreSQL)
	update := sqlbuilder.NewUpdateBuilder()
	updateQ, updateArgs := update.Update(`hr.attendance_periods`).Set(update.Assign(`active`, false)).BuildWithFlavor(sqlbuilder.PostgreSQL)

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

func (repo *PayrollRepository) GetActivePayrollPeriod(ctx context.Context) (PayrollPeriod, error) {
	sq := sqlbuilder.NewSelectBuilder()
	sq.Select(`id`, `start_date`, `end_date`).From(`hr.attendance_periods`).Where(sq.Equal(`active`, true))
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)

	var result PayrollPeriod
	err := tx.QueryRowxContext(ctx, q, args...).StructScan(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (repo *PayrollRepository) IsPayrollCreated(ctx context.Context, periodID string) (bool, error) {
	q := `SELECT 1 FROM hr.payrolls WHERE period_id = $1`

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)

	var result sql.NullInt64
	err := tx.QueryRowxContext(ctx, q, periodID).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
