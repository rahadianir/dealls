package payroll

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rahadianir/dealls/internal/attendance"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/models"
	"github.com/rahadianir/dealls/internal/pkg/dbhelper"
	"github.com/rahadianir/dealls/internal/pkg/xcontext"
	"github.com/rahadianir/dealls/internal/pkg/xerror"
	"github.com/rahadianir/dealls/internal/user"
)

type PayrollLogic struct {
	deps        *config.CommonDependencies
	payrollRepo PayrollRepository
	userRepo    user.UserRepository
	attRepo     attendance.AttendanceRepository
}

func NewPayrollLogic(deps *config.CommonDependencies, payrollRepo PayrollRepository, userRepo user.UserRepository, attRepo attendance.AttendanceRepository) *PayrollLogic {
	return &PayrollLogic{
		deps:        deps,
		payrollRepo: payrollRepo,
		userRepo:    userRepo,
		attRepo:     attRepo,
	}
}

func (logic *PayrollLogic) SetPayrollPeriod(ctx context.Context, start time.Time, end time.Time) error {
	// check admin role of the user
	userID := xcontext.GetUserIDFromContext(ctx)
	isAdmin, err := logic.userRepo.IsAdmin(ctx, userID)
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
	isAdmin, err := logic.userRepo.IsAdmin(ctx, userID)
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

		// get all users reimbursement in the period
		usersReimbursements, err := logic.attRepo.GetAllUserReimbursementsByPeriod(ctx, period.StartDate, period.EndDate)
		if err != nil {
			logic.deps.Logger.ErrorContext(ctx, "failed to get all users reimbursements in payroll period", slog.Any("error", err))
			return err
		}

		// compile all active users in the period
		activeUserMap := make(map[string]PayrollCalculationData)
		activeUserList := []string{}

		for _, att := range usersAttendances {
			_, ok := activeUserMap[att.UserID]
			if !ok {
				activeUserMap[att.UserID] = PayrollCalculationData{
					AttendanceCount: att.Count,
				}
				activeUserList = append(activeUserList, att.UserID)
			}

		}

		for _, ovt := range usersOvertimes {
			activeData, ok := activeUserMap[ovt.UserID]
			if !ok {
				activeUserMap[ovt.UserID] = PayrollCalculationData{
					OvertimeHoursCount: ovt.Count,
				}
				activeUserList = append(activeUserList, ovt.UserID)
			} else {
				activeData.OvertimeHoursCount = ovt.Count
				activeUserMap[ovt.UserID] = activeData
			}
		}

		for _, reimbursement := range usersReimbursements {
			activeData, ok := activeUserMap[reimbursement.UserID]
			if !ok {
				activeUserMap[reimbursement.UserID] = PayrollCalculationData{
					Reimbursements: []models.Reimbursement{reimbursement},
				}
				activeUserList = append(activeUserList, reimbursement.UserID)
			} else {
				activeData.Reimbursements = append(activeData.Reimbursements, reimbursement)
				activeUserMap[reimbursement.UserID] = activeData
			}
		}

		// get all active users salary
		userSalaries, err := logic.userRepo.GetUsersSalaryByIDs(ctx, activeUserList)
		if err != nil {
			logic.deps.Logger.ErrorContext(ctx, "failed to get all active users salaries in payroll period", slog.Any("error", err))
			return err
		}

		// input users salary on calculation data pool
		for _, salary := range userSalaries {
			activeData, ok := activeUserMap[salary.UserID]
			if !ok {
				activeUserMap[salary.UserID] = PayrollCalculationData{
					Salary: salary.Salary,
				}
				activeUserList = append(activeUserList, salary.UserID)
			} else {
				activeData.Salary = salary.Salary
				activeUserMap[salary.UserID] = activeData
			}

		}

		// for id, data := range activeUserMap {
		// 	fmt.Println("userID: ", id)
		// 	fmt.Println("salary: ", data.Salary)
		// 	fmt.Println("attendance count: ", data.AttendanceCount)
		// 	fmt.Println("overtime hour count: ", data.OvertimeHoursCount)
		// 	fmt.Println("reimbursement list: ", data.Reimbursements)
		// }

		return nil
	})
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to process payroll", slog.Any("error", err))
		return err
	}

	return nil
}

func (logic *PayrollLogic) CalculateSalary(ctx context.Context, userID string, attendanceCount int) {}

func (logic *PayrollLogic) CalculateOvertime(ctx context.Context, userID string, overtimeHourCount int) {
}
