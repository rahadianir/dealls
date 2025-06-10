package attendance

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/pkg/xerror"
)

type AttendanceLogic struct {
	deps    *config.CommonDependencies
	attRepo AttendanceRepository
	today   time.Time // putting it here so it's easier to be mocked/tested
}

func NewAttendanceLogic(deps *config.CommonDependencies, attRepo AttendanceRepository) *AttendanceLogic {
	return &AttendanceLogic{
		deps:    deps,
		attRepo: attRepo,
		today:   time.Now(),
	}
}

func (logic *AttendanceLogic) SubmitAttendance(ctx context.Context, userID string, timestamp string) error {
	// parse submitted time
	submittedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		logic.deps.Logger.WarnContext(ctx, "failed to parse timestamp", slog.Any("error", err))
		return xerror.ClientError{Err: err}
	}

	// check whether today is weekend, cuz it's not allowed to submit on weekend
	// check whether submitted day is working hours/day
	today := strings.ToLower(logic.today.Weekday().String())
	submittedTimeDay := strings.ToLower(submittedTime.Weekday().String())
	if today == "saturday" || today == "sunday" || submittedTimeDay == "saturday" || submittedTimeDay == "sunday" {
		return xerror.ClientError{Err: fmt.Errorf("cannot submit attendance in weekend")}
	}

	err = logic.attRepo.SubmitAttendance(ctx, userID, submittedTime)
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to submit attendance", slog.Any("error", err))
		return err
	}

	return nil
}

func (logic *AttendanceLogic) SubmitOvertime(ctx context.Context, userID string, hourCount int, finishedOvertimeTimestamp string) error {
	// parse submitted time
	submittedTime, err := time.Parse(time.RFC3339, finishedOvertimeTimestamp)
	if err != nil {
		logic.deps.Logger.WarnContext(ctx, "failed to parse finished overtime timestamp", slog.Any("error", err))
		return xerror.ClientError{Err: err}
	}

	// check whether overtime is submitted after work hours/day
	submittedDay := submittedTime.Weekday()
	submittedHour := submittedTime.Hour()
	// check submittedDay's day
	if submittedDay >= 1 && submittedDay <= 5 {
		// if overtime is submitted for work days,
		// then check whether the submitted time is already past work time
		if submittedHour >= 9 && submittedHour < 17 {
			return xerror.ClientError{Err: fmt.Errorf("overtime must be submitted outside working hours")}
		}

		// check whether the submitted hours is actual from the last working hours
		// if it's submitted for a different day overtime
		if submittedHour < 9 {
			submittedHour += 24
		}
		if submittedHour-hourCount < 17 {
			return xerror.ClientError{Err: fmt.Errorf("overtime hours overlapped with working hours")}
		}

	}

	// get current overtime hours on that day
	currentOvtHours, err := logic.attRepo.GetUserOvertimeByTime(ctx, userID, submittedTime)
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to get user's overtime hours", slog.Any("error", err))
		return err
	}

	// check whether total overtime hours exceed 3 hours
	if (currentOvtHours + hourCount) > 3 {
		return xerror.ClientError{Err: fmt.Errorf("overtime hours per day cannot exceed 3 hours")}
	}

	err = logic.attRepo.SubmitOvertime(ctx, userID, hourCount, submittedTime)
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to submit overtime hours", slog.Any("error", err))
		return err
	}

	return nil
}

func (logic *AttendanceLogic) SubmitReimbursement(ctx context.Context, userID string, amount float64, desc string) error {
	err := logic.attRepo.SubmitReimbursement(ctx, userID, amount, desc)
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to submit reimbursement", slog.Any("error", err))
		return err
	}
	return nil
}
