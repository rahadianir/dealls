package payroll

import (
	"context"
	"time"

	"github.com/rahadianir/dealls/internal/models"
)

type PayrollRepositoryInterface interface {
	SetPayrollPeriod(ctx context.Context, data PayrollPeriod) error
	GetActivePayrollPeriod(ctx context.Context) (PayrollPeriod, error)
	StorePayslip(ctx context.Context, payslip models.Payslip) error
	MarkPayrollProcessed(ctx context.Context, id string, totalPaid float64) error
	GetPayslipsSummary(ctx context.Context, payrollID string) ([]models.Payslip, error)
	GetUserPayslipByID(ctx context.Context, userID string, payrollID string) (models.Payslip, error)
}

type PayrollLogicInterface interface {
	SetPayrollPeriod(ctx context.Context, start time.Time, end time.Time) error
	CalculatePayroll(ctx context.Context) error
	CalculatePay(ctx context.Context, data PayrollCalculationData) models.Payslip
	GetPayrollsSummary(ctx context.Context) (PayslipSummaryResponse, error)
	GetUserPayslipByID(ctx context.Context, userID string) (models.Payslip, error)
}

