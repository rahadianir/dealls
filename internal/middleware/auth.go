package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/pkg/xcontext"
	"github.com/rahadianir/dealls/internal/pkg/xhttp"
	"github.com/rahadianir/dealls/internal/pkg/xjwt"
)

type AuthMiddleware struct {
	deps      *config.CommonDependencies
	jwtHelper xjwt.JWTHelper
}

func NewAuthMiddleware(deps *config.CommonDependencies, jwtHelper xjwt.JWTHelper) *AuthMiddleware {
	return &AuthMiddleware{
		deps:      deps,
		jwtHelper: jwtHelper,
	}
}

func (mw *AuthMiddleware) AuthOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(http.CanonicalHeaderKey("authorization"))
		if header == "" {
			xhttp.SendJSONResponse(w, xhttp.BaseResponse{
				Error:   "empty header",
				Message: "unauthorized",
			}, http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			xhttp.SendJSONResponse(w, xhttp.BaseResponse{
				Error:   "invalid header",
				Message: "unauthorized",
			}, http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")

		// also calls validate token under
		claims, err := mw.jwtHelper.GetTokenClaims(token, mw.deps.Config.App.JWTSecretKey)
		if err != nil {
			xhttp.SendJSONResponse(w, xhttp.BaseResponse{
				Error:   fmt.Errorf("%s: %w", "invalid token", err).Error(),
				Message: "unauthorized",
			}, http.StatusUnauthorized)
			return
		}

		userID, err := claims.GetSubject()
		if err != nil {
			xhttp.SendJSONResponse(w, xhttp.BaseResponse{
				Error:   fmt.Errorf("%s: %w", "invalid user ID", err).Error(),
				Message: "unauthorized",
			}, http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, xcontext.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
