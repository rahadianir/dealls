package attendance

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/models"
	"github.com/rahadianir/dealls/internal/pkg/dbhelper"
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

func (repo *AttendanceRepository) GetAllUserAttendancesByPeriod(ctx context.Context, start time.Time, end time.Time) ([]models.Attendance, error) {
	sq := sqlbuilder.NewSelectBuilder()
	sq.Select(`count(distinct(user_id, attendance_date)) as count`, `user_id`).From(`hr.attendances`).
		Where(
			sq.And(
				sq.Between(`attendance_date`, start, end),
				sq.IsNull(`deleted_at`),
			),
		).
		GroupBy(`user_id`)
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)

	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil {
		return []models.Attendance{}, err
	}
	defer rows.Close()

	var temp SQLAttendance
	var result []models.Attendance
	for rows.Next() {
		err := rows.StructScan(&temp)
		if err != nil {
			repo.deps.Logger.WarnContext(ctx, "failed to scan attendance data", slog.Any("error", err))
			continue
		}
		result = append(result, models.Attendance{
			UserID: temp.UserID.String,
			Count:  int(temp.Count.Int64),
		})
	}

	return result, nil
}

func (repo *AttendanceRepository) GetAllUserOvertimesByPeriod(ctx context.Context, start time.Time, end time.Time) ([]models.Overtime, error) {
	sq := sqlbuilder.NewSelectBuilder()
	sq.Select(`sum(hour_count) as count`, `user_id`).From(`hr.overtimes`).
		Where(
			sq.And(
				sq.Between(`date`, start, end),
				sq.IsNull(`deleted_at`),
			),
		).
		GroupBy(`user_id`)
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)

	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil {
		return []models.Overtime{}, err
	}
	defer rows.Close()

	var temp SQLOvertime
	var result []models.Overtime
	for rows.Next() {
		err := rows.StructScan(&temp)
		if err != nil {
			repo.deps.Logger.WarnContext(ctx, "failed to scan overtime data", slog.Any("error", err))
			continue
		}
		result = append(result, models.Overtime{
			UserID: temp.UserID.String,
			Count:  int(temp.Count.Int64),
		})
	}

	return result, nil
}

func (repo *AttendanceRepository) GetAllUserReimbursementsByPeriod(ctx context.Context, start time.Time, end time.Time) ([]models.Reimbursement, error) {
	sq := sqlbuilder.NewSelectBuilder()
	sq.Select(`id`, `user_id`, `amount`, `description`).From(`hr.reimbursements`).
		Where(
			sq.And(
				sq.Between(`created_at`, start, end),
				sq.IsNull(`deleted_at`),
			),
		)
	q, args := sq.BuildWithFlavor(sqlbuilder.PostgreSQL)

	tx := dbhelper.ExtractTx(ctx, repo.deps.DB)
	rows, err := tx.QueryxContext(ctx, q, args...)
	if err != nil {
		return []models.Reimbursement{}, err
	}
	defer rows.Close()

	var temp models.SQLReimbursement
	var result []models.Reimbursement
	for rows.Next() {
		err := rows.StructScan(&temp)
		if err != nil {
			repo.deps.Logger.WarnContext(ctx, "failed to scan reimbursement data", slog.Any("error", err))
			continue
		}
		result = append(result, models.Reimbursement{
			ID:          temp.ID.String,
			UserID:      temp.UserID.String,
			Amount:      temp.Amount.Float64,
			Description: temp.Description.String,
		})
	}

	return result, nil
}
