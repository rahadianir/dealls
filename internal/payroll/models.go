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
	ID              string
	StartDate       time.Time
	EndDate         time.Time
	TotalWorkDays   int
	Processed       bool
	TotalSalaryPaid float64
}

type SQLPayrollPeriod struct {
	ID              sql.NullString  `db:"id"`
	StartDate       sql.NullTime    `db:"start_date"`
	EndDate         sql.NullTime    `db:"end_date"`
	TotalWorkDays   sql.NullInt64   `db:"total_work_days"`
	Processed       sql.NullBool    `db:"processed"`
	TotalSalaryPaid sql.NullFloat64 `db:"total_salary_paid"`
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

type SQLPayslip struct {
	ID                 sql.NullString  `db:"id"`
	UserID             sql.NullString  `db:"user_id"`
	TakeHomePay        sql.NullFloat64 `db:"take_home_pay"`
	Name               sql.NullString  `db:"name"`
	PayrollID          sql.NullString  `db:"payroll_id"`
	BaseSalary         sql.NullFloat64 `db:"base_salary"`
	TotalAttendance    sql.NullInt64   `db:"attendance_days"`
	TotalWorkDay       sql.NullInt64   `db:"total_work_days"`
	TotalOvertimeHour  sql.NullInt64   `db:"overtime_hours"`
	OvertimePay        sql.NullFloat64 `db:"overtime_bonus"`
	ReimbursementList  []byte          `db:"reimbursement_list"`
	TotalReimbursement sql.NullFloat64 `db:"total_reimbursement"`
}

type Payslip struct {
	ID                 string
	Name               string
	UserID             string
	PayrollID          string
	BaseSalary         float64
	TotalAttendance    int
	TotalWorkDay       int
	TotalOvertimeHour  int
	OvertimePay        float64
	ReimbursementList  []Reimbursement
	TotalReimbursement float64
	TakeHomePay        float64
}

type PayslipResponse struct {
	UserID      string  `json:"user_id"`
	TakeHomePay float64 `json:"take_home_pay"`
	Name        string  `json:"name"`
}
type PayslipSummaryResponse struct {
	TotalTakeHomePay float64           `json:"total_take_home_pay"`
	Payslips         []PayslipResponse `json:"payslips"`
}

type UserPayslipRequest struct {
	UserID string `json:"user_id"`
}
