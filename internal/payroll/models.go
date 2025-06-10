package payroll

import (
	"time"

	"github.com/rahadianir/dealls/internal/models"
)

type PayrollPeriodRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type PayrollPeriod struct {
	ID        string    `db:"id"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
}

type PayrollCalculationData struct {
	AttendanceCount    int
	OvertimeHoursCount int
	Reimbursements     []models.Reimbursement
	Salary             float64
}
