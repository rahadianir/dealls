package models

import "database/sql"

type SQLReimbursement struct {
	ID          sql.NullString
	UserID      sql.NullString `db:"user_id"`
	Amount      sql.NullFloat64
	Description sql.NullString
}

type Reimbursement struct {
	ID          string  `json:"id,omitempty"`
	UserID      string  `json:"user_id,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	Description string  `json:"description,omitempty"`
}

type Payslip struct {
	ID                 string          `json:"id"`
	Name               string          `json:"name"`
	UserID             string          `json:"user_id"`
	PayrollID          string          `json:"payroll_id"`
	BaseSalary         float64         `json:"base_salary"`
	TotalAttendance    int             `json:"total_attendance"`
	TotalWorkDay       int             `json:"total_work_day"`
	TotalOvertimeHour  int             `json:"total_overtime_hour"`
	OvertimePay        float64         `json:"overtime_bonus"`
	ReimbursementList  []Reimbursement `json:"reimbursement_list"`
	TotalReimbursement float64         `json:"total_reimbursement_amount"`
	TakeHomePay        float64         `json:"take_home_pay"`
}
