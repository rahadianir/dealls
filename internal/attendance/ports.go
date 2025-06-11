package attendance

import (
	"context"
	"time"

	"github.com/rahadianir/dealls/internal/models"
)

type AttendanceRepositoryInterface interface {
	SubmitAttendance(ctx context.Context, userID string, timestamp time.Time) error
	SubmitOvertime(ctx context.Context, userID string, hours int, timestamp time.Time) error
	GetUserOvertimeByTime(ctx context.Context, userID string, date time.Time) (int, error)
	SubmitReimbursement(ctx context.Context, userID string, amount float64, desc string) error
	GetAllUserAttendancesByPeriod(ctx context.Context, start time.Time, end time.Time) ([]models.Attendance, error)
	GetAllUserOvertimesByPeriod(ctx context.Context, start time.Time, end time.Time) ([]models.Overtime, error)
	GetAllUserReimbursementsByPeriod(ctx context.Context, start time.Time, end time.Time) ([]models.Reimbursement, error)
}

type AttendanceLogicInterface interface {
	SubmitAttendance(ctx context.Context, userID string, timestamp string) error
	SubmitOvertime(ctx context.Context, userID string, hourCount int, finishedOvertimeTimestamp string) error
	SubmitReimbursement(ctx context.Context, userID string, amount float64, desc string) error
}
