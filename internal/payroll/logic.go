package payroll

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
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

	totalWorkDay := calculateWorkingDays(start, end)

	err = logic.payrollRepo.SetPayrollPeriod(ctx, PayrollPeriod{
		ID:            uuid.NewString(),
		StartDate:     start,
		EndDate:       end,
		TotalWorkDays: totalWorkDay,
	})
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

		// check whether payroll is already processed
		if period.Processed {
			return xerror.LogicError{Err: fmt.Errorf("payroll processed already!")}
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

		// compile all active users and other related data in the period

		// setup map to store all payroll related data
		activeUserMap := make(map[string]PayrollCalculationData)

		// setup list to store all active user IDs (user that worked in the active payroll period)
		activeUserList := []string{}

		// populate payroll and active user data with attendance data
		for _, att := range usersAttendances {
			_, ok := activeUserMap[att.UserID]
			if !ok {
				activeUserMap[att.UserID] = PayrollCalculationData{
					UserID:          att.UserID,
					PayrollID:       period.ID,
					TotalWorkDay:    period.TotalWorkDays,
					AttendanceCount: att.Count,
				}
				activeUserList = append(activeUserList, att.UserID)
			}

		}

		// populate payroll and active user data with overtime data
		for _, ovt := range usersOvertimes {
			activeData, ok := activeUserMap[ovt.UserID]
			if !ok {
				activeUserMap[ovt.UserID] = PayrollCalculationData{
					UserID:             ovt.UserID,
					PayrollID:          period.ID,
					TotalWorkDay:       period.TotalWorkDays,
					OvertimeHoursCount: ovt.Count,
				}
				activeUserList = append(activeUserList, ovt.UserID)
			} else {
				activeData.OvertimeHoursCount = ovt.Count
				activeUserMap[ovt.UserID] = activeData
			}
		}

		// populate payroll and active user data with reimbursement data
		for _, reimbursement := range usersReimbursements {
			activeData, ok := activeUserMap[reimbursement.UserID]
			if !ok {
				activeUserMap[reimbursement.UserID] = PayrollCalculationData{
					UserID:       reimbursement.UserID,
					PayrollID:    period.ID,
					TotalWorkDay: period.TotalWorkDays,
					Reimbursements: []Reimbursement{
						{
							ID:     reimbursement.ID,
							Amount: reimbursement.Amount,
							Desc:   reimbursement.Description,
						},
					},
				}
				activeUserList = append(activeUserList, reimbursement.UserID)
			} else {
				activeData.Reimbursements = append(activeData.Reimbursements, Reimbursement{
					ID:     reimbursement.ID,
					Amount: reimbursement.Amount,
					Desc:   reimbursement.Description,
				})
				activeUserMap[reimbursement.UserID] = activeData
			}
		}

		// get all active users salary
		userSalaries, err := logic.userRepo.GetUsersSalaryByIDs(ctx, activeUserList)
		if err != nil {
			logic.deps.Logger.ErrorContext(ctx, "failed to get all active users salaries in payroll period", slog.Any("error", err))
			return err
		}

		// populate payroll and active user data with salary data
		for _, salary := range userSalaries {
			activeData, ok := activeUserMap[salary.UserID]
			if !ok {
				activeUserMap[salary.UserID] = PayrollCalculationData{
					UserID:       salary.UserID,
					PayrollID:    period.ID,
					TotalWorkDay: period.TotalWorkDays,
					Salary:       salary.Salary,
				}

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

		// setup worker to calculate payroll
		// setup wait group for flow control
		var wg sync.WaitGroup

		// setup channel to pass the calculation data and result
		jobChan := make(chan PayrollCalculationData)
		payslipChan := make(chan models.Payslip)

		// spawn worker to consume data and process calculation
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for job := range jobChan {
					payslipChan <- logic.CalculatePay(ctx, job)
				}
			}()
		}

		// feed data through channel
		go func() {
			for _, data := range activeUserMap {
				jobChan <- data
			}
			close(jobChan)
		}()

		// setup total salary paid
		var totalSalaryPaid float64

		// spawn worker that receives calculation result
		// and store it to database
		go func() {
			for payslip := range payslipChan {
				totalSalaryPaid += payslip.TakeHomePay
				err := logic.payrollRepo.StorePayslip(ctx, payslip)
				if err != nil {
					logic.deps.Logger.ErrorContext(ctx, "failed to store payslip data", slog.Any("error", err))
					return
				}
			}
		}()

		wg.Wait()
		close(payslipChan)

		err = logic.payrollRepo.MarkPayrollProcessed(ctx, period.ID, totalSalaryPaid)
		if err != nil {
			logic.deps.Logger.ErrorContext(ctx, "failed to mark payroll period processed", slog.Any("error", err))
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

func (logic *PayrollLogic) CalculatePay(ctx context.Context, data PayrollCalculationData) models.Payslip {
	payslip := models.Payslip{
		ID:                uuid.NewString(),
		UserID:            data.UserID,
		PayrollID:         data.PayrollID,
		BaseSalary:        data.Salary,
		TotalAttendance:   data.AttendanceCount,
		TotalWorkDay:      data.TotalWorkDay,
		TotalOvertimeHour: data.OvertimeHoursCount,
	}

	// calculate prorated salary = (total attendance / total work day) * salary
	salary := (float64(payslip.TotalAttendance) / float64(payslip.TotalWorkDay)) * (payslip.BaseSalary)

	// calculate overtime pay = prorated salary per hour * overtime hour
	overtime := (payslip.BaseSalary / float64(payslip.TotalWorkDay) / 8) * float64(payslip.TotalOvertimeHour)
	payslip.OvertimePay = overtime

	// calculate reimbursement
	var reimburseAmount float64
	for _, r := range data.Reimbursements {
		reimburseAmount += r.Amount
		payslip.ReimbursementList = append(payslip.ReimbursementList, models.Reimbursement{
			ID:          r.ID,
			Amount:      r.Amount,
			Description: r.Desc,
		})
	}
	payslip.TotalReimbursement = reimburseAmount

	payslip.TakeHomePay = salary + overtime + reimburseAmount

	return payslip
}

func calculateWorkingDays(startTime time.Time, endTime time.Time) int {
	// Reduce dates to previous Mondays
	startOffset := weekday(startTime)
	startTime = startTime.AddDate(0, 0, -startOffset)
	endOffset := weekday(endTime)
	endTime = endTime.AddDate(0, 0, -endOffset)

	// Calculate weeks and days
	dif := endTime.Sub(startTime)
	weeks := int(math.Round((dif.Hours() / 24) / 7))
	days := -min(startOffset, 5) + min(endOffset, 5)

	// Calculate total days
	return weeks*5 + days
}

func weekday(d time.Time) int {
	wd := d.Weekday()
	if wd == time.Sunday {
		return 6
	}
	return int(wd) - 1
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
