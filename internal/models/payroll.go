package models

import "database/sql"

type SQLReimbursement struct {
	ID          sql.NullString
	UserID      sql.NullString `db:"user_id"`
	Amount      sql.NullFloat64
	Description sql.NullString
}

type Reimbursement struct {
	ID          string
	UserID      string
	Amount      float64
	Description string
}

type Payslip struct {
	ID                 string
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
