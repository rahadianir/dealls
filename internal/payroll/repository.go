package payroll

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/huandu/go-sqlbuilder"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/models"
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

func (repo *PayrollRepository) SetPayrollPeriod(ctx context.Context, data PayrollPeriod) error {
	ins := sqlbuilder.NewInsertBuilder()
	insertQ, insertArgs := ins.InsertInto(`hr.payrolls`).
		Cols(`id`, `start_date`, `end_date`, `active`, `created_at`, `total_work_days`).
		Values(data.ID, data.StartDate, data.EndDate, true, `now()`, data.TotalWorkDays).BuildWithFlavor(sqlbuilder.PostgreSQL)
	update := sqlbuilder.NewUpdateBuilder()
	updateQ, updateArgs := update.Update(`hr.payrolls`).Set(update.Assign(`active`, nil)).BuildWithFlavor(sqlbuilder.PostgreSQL)

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
	sq.Select(`id`, `start_date`, `end_date`, `total_work_days`, `processed`).From(`hr.payrolls`).Where(sq.Equal(`active`, true))
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)

	var temp SQLPayrollPeriod
	err := tx.QueryRowxContext(ctx, q, args...).StructScan(&temp)
	if err != nil {
		return PayrollPeriod{}, err
	}

	result := PayrollPeriod{
		ID:            temp.ID.String,
		StartDate:     temp.StartDate.Time,
		EndDate:       temp.EndDate.Time,
		TotalWorkDays: int(temp.TotalWorkDays.Int64),
		Processed:     temp.Processed.Bool,
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

func (repo *PayrollRepository) StorePayslip(ctx context.Context, payslip models.Payslip) error {
	// convert reimbursement list to JSON first
	reimbursementList := `{}`
	if len(payslip.ReimbursementList) != 0 {
		dataBytes, err := json.Marshal(payslip.ReimbursementList)
		if err != nil {
			repo.deps.Logger.ErrorContext(ctx, "failed to marshal reimbursement list to payslip", slog.Any("error", err))
			return err
		}
		reimbursementList = string(dataBytes)
	}

	sq := sqlbuilder.NewInsertBuilder()
	sq.InsertInto(`hr.payslips`).
		Cols(`id`, `payroll_id`, `user_id`, `base_salary`, `attendance_days`, `total_work_days`, `overtime_hours`, `overtime_bonus`, `reimbursement_list`, `total_reimbursement`, `take_home_pay`, `created_at`).
		Values(payslip.ID, payslip.PayrollID, payslip.UserID, payslip.BaseSalary, payslip.TotalAttendance, payslip.TotalWorkDay, payslip.TotalOvertimeHour, payslip.OvertimePay, reimbursementList, payslip.TotalReimbursement, payslip.TakeHomePay, `now()`)

	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)

	_, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PayrollRepository) MarkPayrollProcessed(ctx context.Context, id string, totalPaid float64) error {
	sq := sqlbuilder.NewUpdateBuilder()
	sq.Update(`hr.payrolls`).Set(
		sq.Assign(`processed`, true),
		sq.Assign(`total_salary_paid`, totalPaid),
	).Where(
		sq.EQ(`id`, id),
	)
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)

	_, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}
