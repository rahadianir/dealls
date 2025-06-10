package payroll

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/huandu/go-sqlbuilder"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/models"
	"github.com/rahadianir/dealls/internal/pkg/dbhelper"
	"github.com/rahadianir/dealls/internal/pkg/xerror"
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
	sq.Select(`id`, `start_date`, `end_date`, `total_work_days`, `processed`, `total_salary_paid`).From(`hr.payrolls`).Where(sq.Equal(`active`, true))
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)

	var temp SQLPayrollPeriod
	err := tx.QueryRowxContext(ctx, q, args...).StructScan(&temp)
	if err != nil {
		return PayrollPeriod{}, err
	}

	result := PayrollPeriod{
		ID:              temp.ID.String,
		StartDate:       temp.StartDate.Time,
		EndDate:         temp.EndDate.Time,
		TotalWorkDays:   int(temp.TotalWorkDays.Int64),
		Processed:       temp.Processed.Bool,
		TotalSalaryPaid: temp.TotalSalaryPaid.Float64,
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

func (repo *PayrollRepository) GetPayslipsSummary(ctx context.Context, payrollID string) ([]models.Payslip, error) {
	sq := sqlbuilder.NewSelectBuilder()
	sq.Select(`p.user_id`, `p.take_home_pay`, `u.name`).From(`hr.payslips p `).Join(`hr.users u`, `p.user_id = u.id`).Where(sq.Equal(`p.payroll_id`, payrollID))
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)

	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil {
		return []models.Payslip{}, err
	}
	defer rows.Close()

	var temp SQLPayslip
	var result []models.Payslip
	for rows.Next() {
		err := rows.StructScan(&temp)
		if err != nil {
			repo.deps.Logger.WarnContext(ctx, "failed to scan payslip summary", slog.Any("error", err))
			continue
		}

		result = append(result, models.Payslip{
			UserID:      temp.UserID.String,
			Name:        temp.Name.String,
			TakeHomePay: temp.TakeHomePay.Float64,
		})
	}

	return result, nil
}

func (repo *PayrollRepository) GetUserPayslipByID(ctx context.Context, userID string, payrollID string) (models.Payslip, error) {
	sq := sqlbuilder.NewSelectBuilder()
	sq.Select(`p.id`, `u.name`, `user_id`, `p.base_salary`, `attendance_days`, `total_work_days`, `overtime_hours`, `overtime_bonus`, `reimbursement_list`, `total_reimbursement`, `take_home_pay`).
		From(`hr.payslips p`).Join(`hr.users u`, `p.user_id = u.id`).Where(
		sq.And(
			sq.Equal(`user_id`, userID),
			sq.Equal(`payroll_id`, payrollID),
		),
	)
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)

	var temp SQLPayslip
	err := tx.QueryRowxContext(ctx, q, args...).StructScan(&temp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Payslip{}, xerror.ErrDataNotFound
		}

		return models.Payslip{}, err
	}

	var list []models.Reimbursement
	if len(temp.ReimbursementList) != 0 {
		err := json.Unmarshal(temp.ReimbursementList, &list)
		if err != nil {
			return models.Payslip{}, fmt.Errorf("failed to unmarshal reimbursement list: %w", err)
		}
	}

	result := models.Payslip{
		ID:                 temp.ID.String,
		Name:               temp.Name.String,
		UserID:             temp.UserID.String,
		PayrollID:          temp.PayrollID.String,
		BaseSalary:         temp.BaseSalary.Float64,
		TotalAttendance:    int(temp.TotalAttendance.Int64),
		TotalWorkDay:       int(temp.TotalWorkDay.Int64),
		TotalOvertimeHour:  int(temp.TotalOvertimeHour.Int64),
		OvertimePay:        temp.OvertimePay.Float64,
		ReimbursementList:  list,
		TotalReimbursement: temp.TotalReimbursement.Float64,
		TakeHomePay:        temp.TakeHomePay.Float64,
	}

	return result, nil
}
