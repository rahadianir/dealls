package attendance

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/rahadianir/dealls/internal/config"
)

type AttendanceRepository struct {
	deps *config.CommonDependencies
}

func NewAttendanceRepository(deps *config.CommonDependencies) *AttendanceRepository {
	return &AttendanceRepository{
		deps: deps,
	}
}

func (repo *AttendanceRepository) SubmitAttendance(ctx context.Context, userID string, timestamp time.Time) error {
	sq := sqlbuilder.NewInsertBuilder()
	q, args := sq.InsertInto(`hr.attendances`).
		Cols(`id`, `user_id`, `attendance_time`, `attendance_date`, `created_at`, `created_by`).
		Values(uuid.NewString(), userID, timestamp, timestamp, `now()`, userID).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx, err := repo.deps.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repo *AttendanceRepository) SubmitOvertime(ctx context.Context, userID string, hours int) error {
	sq := sqlbuilder.NewInsertBuilder()
	q, args := sq.InsertInto(`hr.overtimes`).
		Cols(`id`, `user_id`, `date`, `hour_count`, `created_at`, `created_by`).
		Values(uuid.NewString(), userID, `now()`, hours, `now()`, userID).
		BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx, err := repo.deps.DB.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
	return nil
}

func (repo *AttendanceRepository) SubmitReimbursement(ctx context.Context, userID string, amount float64, desc string) error {
	return nil
}
