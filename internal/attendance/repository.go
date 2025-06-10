package attendance

import (
	"context"
	"database/sql"
	"errors"
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

func (repo *AttendanceRepository) SubmitOvertime(ctx context.Context, userID string, hours int, timestamp time.Time) error {
	sq := sqlbuilder.NewInsertBuilder()
	q, args := sq.InsertInto(`hr.overtimes`).
		Cols(`id`, `user_id`, `date`, `hour_count`, `created_at`, `created_by`).
		Values(uuid.NewString(), userID, timestamp, hours, `now()`, userID).
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

func (repo *AttendanceRepository) GetUserOvertimeByTime(ctx context.Context, userID string, date time.Time) (int, error) {
	sq := sqlbuilder.NewSelectBuilder()
	sq.Select(`SUM(hour_count)`).From(`hr.overtimes`).Where(
		sq.And(
			sq.Equal(`date`, date),
			sq.Equal(`user_id`, userID),
			sq.IsNull(`deleted_at`),
		),
	)
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var overtimeHours sql.NullInt64
	err := repo.deps.DB.QueryRowxContext(ctx, q, args...).Scan(&overtimeHours)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return int(overtimeHours.Int64), nil

}

func (repo *AttendanceRepository) SubmitReimbursement(ctx context.Context, userID string, amount float64, desc string) error {
	sq := sqlbuilder.NewInsertBuilder()
	q, args := sq.InsertInto(`hr.reimbursements`).
		Cols(`id`, `user_id`, `amount`, `description`, `created_at`, `created_by`).
		Values(uuid.NewString(), userID, amount, desc, `now()`, userID).
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
