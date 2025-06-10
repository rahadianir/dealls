package payroll

import "time"

type PayrollPeriodRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type PayrollPeriod struct {
	ID        string    `db:"id"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
}
