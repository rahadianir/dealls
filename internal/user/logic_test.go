package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/models"
	"github.com/rahadianir/dealls/internal/pkg/xjwt"
	"go.uber.org/mock/gomock"
)

func TestUserLogic_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockUserRepositoryInterface(ctrl)
	mockJwt := xjwt.NewMockJWTHelper(ctrl)
	mockDeps := config.CommonDependencies{
		Config: config.InitConfig(context.Background()),
	}

	type fields struct {
		deps      *config.CommonDependencies
		userRepo  UserRepositoryInterface
		jwtHelper xjwt.JWTHelper
	}
	type args struct {
		ctx      context.Context
		username string
		password string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      LoginResponse
		wantErr   bool
		behaviour func()
	}{
		// TODO: Add test cases.
		{
			name: "success login",
			fields: fields{
				deps:      &mockDeps,
				userRepo:  mockRepo,
				jwtHelper: mockJwt,
			},
			args: args{
				ctx:      context.Background(),
				username: "admin",
				password: "admin",
			},
			want: LoginResponse{
				Token: "token",
			},
			wantErr: false,
			behaviour: func() {
				mockRepo.EXPECT().GetUserDetailsByUsername(gomock.Any(), "admin").Return(models.User{
					ID:       "1",
					Password: "$2a$12$x57I28hfnEEJGXE5splrqeNLwWSlhXyFaoDZamMJc9oElJgpUPbwe", // hashed "admin"
				}, nil)
				mockJwt.EXPECT().GenerateJWT(gomock.Any(), "1", gomock.Any(), gomock.Any()).Return("token", nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := &UserLogic{
				deps:      tt.fields.deps,
				userRepo:  tt.fields.userRepo,
				jwtHelper: tt.fields.jwtHelper,
			}
			tt.behaviour()
			got, err := logic.Login(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserLogic.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserLogic.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
