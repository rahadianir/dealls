package payroll

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rahadianir/dealls/internal/attendance"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/pkg/dbhelper"
	"github.com/rahadianir/dealls/internal/pkg/xcontext"
	"github.com/rahadianir/dealls/internal/pkg/xerror"
	"github.com/rahadianir/dealls/internal/user"
)

type PayrollLogic struct {
	deps        *config.CommonDependencies
	payrollRepo PayrollRepository
	userLogic   user.UserLogic
	attRepo     attendance.AttendanceRepository
}

func NewPayrollLogic(deps *config.CommonDependencies, payrollRepo PayrollRepository, userLogic user.UserLogic, attRepo attendance.AttendanceRepository) *PayrollLogic {
	return &PayrollLogic{
		deps:        deps,
		payrollRepo: payrollRepo,
		userLogic:   userLogic,
		attRepo:     attRepo,
	}
}

func (logic *PayrollLogic) SetPayrollPeriod(ctx context.Context, start time.Time, end time.Time) error {
	// check admin role of the user
	userID := xcontext.GetUserIDFromContext(ctx)
	isAdmin, err := logic.userLogic.IsAdmin(ctx, userID)
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to check user admin role", slog.Any("error", err))
		return err
	}

	if !isAdmin {
		return xerror.AuthError{Err: fmt.Errorf("admin only operation")}
	}

	err = logic.payrollRepo.SetPayrollPeriod(ctx, start, end)
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to set payroll period", slog.Any("error", err))
		return err
	}

	return nil
}

func (logic *PayrollLogic) CalculatePayroll(ctx context.Context) error {
	// check admin role of the user
	userID := xcontext.GetUserIDFromContext(ctx)
	isAdmin, err := logic.userLogic.IsAdmin(ctx, userID)
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to check user admin role", slog.Any("error", err))
		return err
	}

	if !isAdmin {
		return xerror.AuthError{Err: fmt.Errorf("admin only operation")}
	}

	err = dbhelper.WithTransaction(ctx, logic.deps.DB, func(ctx context.Context) error {
		// get active payroll period
		period, err := logic.payrollRepo.GetActivePayrollPeriod(ctx)
		if err != nil {
			logic.deps.Logger.ErrorContext(ctx, "failed to get active payroll period", slog.Any("error", err))
			return err
		}

		// check whether payroll is already generated
		isGenerated, err := logic.payrollRepo.IsPayrollCreated(ctx, period.ID)
		if err != nil {
			logic.deps.Logger.ErrorContext(ctx, "failed to check whether payroll is created already", slog.Any("error", err))
			return err
		}

		if isGenerated {
			return xerror.LogicError{Err: fmt.Errorf("payroll generated already!")}
		}

		// get all users attendances in the period
		usersAttendances, err := logic.attRepo.GetAllUserAttendancesByPeriod(ctx, period.StartDate, period.EndDate)
		if err != nil {
			logic.deps.Logger.ErrorContext(ctx, "failed to get all users attendances in payroll period", slog.Any("error", err))
			return err
		}

		// get all users overtimes in the period
		usersOvertimes, err := logic.attRepo.GetAllUserOvertimesByPeriod(ctx, period.StartDate, period.EndDate)
		if err != nil {
			logic.deps.Logger.ErrorContext(ctx, "failed to get all users overtimes in payroll period", slog.Any("error", err))
			return err
		}

		return nil
	})
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to process payroll", slog.Any("error", err))
		return err
	}

	return nil
}
