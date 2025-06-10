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
	}, http.StatusCreated)
}
