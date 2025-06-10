package payroll

import (
	"database/sql"
	"time"
)

type PayrollPeriodRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type PayrollPeriod struct {
	ID            string
	StartDate     time.Time
	EndDate       time.Time
	TotalWorkDays int
}

type SQLPayrollPeriod struct {
	ID            sql.NullString `db:"id"`
	StartDate     sql.NullTime   `db:"start_date"`
	EndDate       sql.NullTime   `db:"end_date"`
	TotalWorkDays sql.NullInt64  `db:"total_work_days"`
}

type Reimbursement struct {
	ID     string
	Amount float64
	Desc   string
}
type PayrollCalculationData struct {
	UserID             string
	PayrollID          string
	TotalWorkDay       int
	AttendanceCount    int
	OvertimeHoursCount int
	Reimbursements     []Reimbursement
	Salary             float64
}
