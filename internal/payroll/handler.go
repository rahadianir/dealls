package payroll

import (
	"net/http"

	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/pkg/xerror"
	"github.com/rahadianir/dealls/internal/pkg/xhttp"
)

type PayrollHandler struct {
	deps         *config.CommonDependencies
	payrollLogic PayrollLogic
}

func NewPayrollHandler(deps *config.CommonDependencies, payrollLogic PayrollLogic) *PayrollHandler {
	return &PayrollHandler{
		deps:         deps,
		payrollLogic: payrollLogic,
	}
}

func (h *PayrollHandler) SetPayrollPeriod(w http.ResponseWriter, r *http.Request) {
	var payload PayrollPeriodRequest
	err := xhttp.BindJSONRequest(r, &payload)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: xerror.ErrBadRequest.Error(),
		}, http.StatusBadRequest)
		return
	}

	err = h.payrollLogic.SetPayrollPeriod(r.Context(), payload.StartDate, payload.EndDate)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to set payroll period",
		}, http.StatusBadRequest)
		return
	}

	xhttp.SendJSONResponse(w, xhttp.BaseResponse{
		Message: "payroll period set",
	}, http.StatusOK)
}

func (h *PayrollHandler) CalculatePayroll(w http.ResponseWriter, r *http.Request) {
	err := h.payrollLogic.CalculatePayroll(r.Context())
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to calculate payroll in active period",
		}, http.StatusBadRequest)
		return
	}

	xhttp.SendJSONResponse(w, xhttp.BaseResponse{
		Message: "payroll in active period calculated",
	}, http.StatusOK)
}

func (h *PayrollHandler) GeneratePayrollSummary(w http.ResponseWriter, r *http.Request) {
	resp, err := h.payrollLogic.GetPayrollsSummary(r.Context())
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to get payroll summary in active period",
		}, http.StatusBadRequest)
		return
	}
	xhttp.SendJSONResponse(w, xhttp.BaseResponse{
		Message: "payroll summary in active period generated",
		Data:    resp,
	}, http.StatusOK)
}

func (h *PayrollHandler) GetUserPayslip(w http.ResponseWriter, r *http.Request) {
	var payload UserPayslipRequest
	err := xhttp.BindJSONRequest(r, &payload)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: xerror.ErrBadRequest.Error(),
		}, http.StatusBadRequest)
		return
	}

	data, err := h.payrollLogic.GetUserPayslipByID(r.Context(), payload.UserID)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to get user payslip summary in active period",
		}, http.StatusBadRequest)
		return
	}
	xhttp.SendJSONResponse(w, xhttp.BaseResponse{
		Message: "payslip summary in active period fetched",
		Data:    data,
	}, http.StatusOK)
}
