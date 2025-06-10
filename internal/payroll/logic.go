package payroll

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/pkg/xcontext"
	"github.com/rahadianir/dealls/internal/pkg/xerror"
	"github.com/rahadianir/dealls/internal/user"
)

type PayrollLogic struct {
	deps        *config.CommonDependencies
	payrollRepo PayrollRepository
	userLogic   user.UserLogic
}

func NewPayrollLogic(deps *config.CommonDependencies, payrollRepo PayrollRepository, userLogic user.UserLogic) *PayrollLogic {
	return &PayrollLogic{
		deps:        deps,
		payrollRepo: payrollRepo,
		userLogic:   userLogic,
	}
}

func (logic *PayrollLogic) SetPayrollPeriod(ctx context.Context, start time.Time, end time.Time) error {
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
