package attendance

import "database/sql"

type AttendanceRequest struct {
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
}

type OvertimeRequest struct {
	UserID    string `json:"user_id"`
	Hours     int    `json:"hours"`
	Timestamp string `json:"timestamp"`
}

type ReimbursementRequest struct {
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type SQLAttendance struct {
	UserID sql.NullString `db:"user_id"`
	Count  sql.NullInt64  `db:"count"`
}

type SQLOvertime struct {
	UserID sql.NullString `db:"user_id"`
	Count  sql.NullInt64  `db:"count"`
}
