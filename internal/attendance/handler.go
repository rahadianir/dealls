package attendance

import (
	"net/http"

	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/pkg/xerror"
	"github.com/rahadianir/dealls/internal/pkg/xhttp"
)

type AttendanceHandler struct {
	deps     *config.CommonDependencies
	attLogic AttendanceLogic
}

func NewAttendanceHandler(deps *config.CommonDependencies, attLogic AttendanceLogic) *AttendanceHandler {
	return &AttendanceHandler{
		deps:     deps,
		attLogic: attLogic,
	}
}

func (h *AttendanceHandler) SubmitAttendance(w http.ResponseWriter, r *http.Request) {
	var payload AttendanceRequest
	err := xhttp.BindJSONRequest(r, &payload)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: xerror.ErrBadRequest.Error(),
		}, http.StatusBadRequest)
		return
	}

	err = h.attLogic.SubmitAttendance(r.Context(), payload.UserID, payload.Timestamp)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to submit attendance",
		}, http.StatusBadRequest)
		return
	}

	xhttp.SendJSONResponse(w, xhttp.BaseResponse{
		Message: "attendance submitted",
	}, http.StatusCreated)
}

func (h *AttendanceHandler) SubmitOvertime(w http.ResponseWriter, r *http.Request) {
	var payload OvertimeRequest
	err := xhttp.BindJSONRequest(r, &payload)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: xerror.ErrBadRequest.Error(),
		}, http.StatusBadRequest)
		return
	}

	err = h.attLogic.SubmitOvertime(r.Context(), payload.UserID, payload.Hours, payload.Timestamp)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to submit overtime",
		}, http.StatusBadRequest)
		return
	}

	xhttp.SendJSONResponse(w, xhttp.BaseResponse{
		Message: "overtime submitted",
	}, http.StatusCreated)

}

func (h *AttendanceHandler) SubmitReimbursement(w http.ResponseWriter, r *http.Request) {
	var payload ReimbursementRequest
	err := xhttp.BindJSONRequest(r, &payload)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: xerror.ErrBadRequest.Error(),
		}, http.StatusBadRequest)
		return
	}

	err = h.attLogic.SubmitReimbursement(r.Context(), payload.UserID, payload.Amount, payload.Description)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to submit reimbursement",
		}, http.StatusBadRequest)
		return
	}

	xhttp.SendJSONResponse(w, xhttp.BaseResponse{
		Message: "reimbursement submitted",
	}, http.StatusCreated)
}
