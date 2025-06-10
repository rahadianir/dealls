package payroll

import "time"

type PayrollPeriodRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
