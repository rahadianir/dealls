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
}

func NewAttendanceLogic(deps *config.CommonDependencies, attRepo AttendanceRepository) *AttendanceLogic {
	return &AttendanceLogic{
		deps:    deps,
		attRepo: attRepo,
	}
}

func (logic *AttendanceLogic) SubmitAttendance(ctx context.Context, userID string) error {
	today := strings.ToLower(time.Now().Weekday().String())
	if today == "saturday" || today == "sunday" {
		return xerror.ClientError{Err: fmt.Errorf("cannot submit attendance in weekend")}
	}

	err := logic.attRepo.SubmitAttendance(ctx, userID, time.Now())
	if err != nil {
		logic.deps.Logger.ErrorContext(ctx, "failed to submit attendance", slog.Any("error", err))
		return err
	}

	return nil
}

func (logic *AttendanceLogic) SubmitOvertime(ctx context.Context, hours int) error {
	return nil
}

func (logic *AttendanceLogic) SubmitReimbursement(ctx context.Context, userID string, amount float64, desc string) error {
	return nil
}
