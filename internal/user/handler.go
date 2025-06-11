package user

import (
	"net/http"

	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/pkg/xerror"
	"github.com/rahadianir/dealls/internal/pkg/xhttp"
)

type UserHandler struct {
	deps      *config.CommonDependencies
	userLogic UserLogicInterface
}

func NewUserHandler(deps *config.CommonDependencies, userLogic UserLogicInterface) *UserHandler {
	return &UserHandler{
		deps:      deps,
		userLogic: userLogic,
	}
}

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload LoginRequest
	err := xhttp.BindJSONRequest(r, &payload)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: xerror.ErrBadRequest.Error(),
		}, http.StatusBadRequest)
		return
	}

	result, err := handler.userLogic.Login(r.Context(), payload.Username, payload.Password)
	if err != nil {
		xhttp.SendJSONResponse(w, xhttp.BaseResponse{
			Error:   err.Error(),
			Message: "failed to login",
		}, xerror.ParseErrorTypeToCodeInt(err))
		return
	}

	xhttp.SendJSONResponse(w, xhttp.BaseResponse{
		Message: "login success",
		Data:    result,
	}, http.StatusOK)
}
