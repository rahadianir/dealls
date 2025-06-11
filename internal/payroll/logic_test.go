package payroll

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/rahadianir/dealls/internal/attendance"
	"github.com/rahadianir/dealls/internal/config"
	"github.com/rahadianir/dealls/internal/pkg/xcontext"
	"github.com/rahadianir/dealls/internal/user"
	"go.uber.org/mock/gomock"
)

func TestPayrollLogic_SetPayrollPeriod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDeps := config.CommonDependencies{
		Config: config.InitConfig(context.Background()),
	}

	mockPayrollRepo := NewMockPayrollRepositoryInterface(ctrl)
	mockUserRepo := user.NewMockUserRepositoryInterface(ctrl)
	mockAttRepo := attendance.NewMockAttendanceRepositoryInterface(ctrl)
	startTime, err := time.Parse(time.RFC3339, "2025-05-25T06:29:44+07:00")
	if err != nil {
		t.Fatal(err)
	}
	endTime, err := time.Parse(time.RFC3339, "2025-06-25T06:29:44+07:00")
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		deps        *config.CommonDependencies
		payrollRepo PayrollRepositoryInterface
		userRepo    user.UserRepositoryInterface
		attRepo     attendance.AttendanceRepositoryInterface
	}
	type args struct {
		ctx   context.Context
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		behaviour func(f fields, a args)
	}{
		// TODO: Add test cases.
		{
			name: "success set payroll period",
			fields: fields{
				deps:        &mockDeps,
				payrollRepo: mockPayrollRepo,
				userRepo:    mockUserRepo,
				attRepo:     mockAttRepo,
			},
			args: args{
				ctx:   context.WithValue(context.Background(), xcontext.UserIDKey, "user-id"),
				start: startTime,
				end:   endTime,
			},
			wantErr: false,
			behaviour: func(f fields, a args) {
				mockUserRepo.EXPECT().IsAdmin(gomock.Any(), "user-id").Return(true, nil)
				mockPayrollRepo.EXPECT().SetPayrollPeriod(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := &PayrollLogic{
				deps:        tt.fields.deps,
				payrollRepo: tt.fields.payrollRepo,
				userRepo:    tt.fields.userRepo,
				attRepo:     tt.fields.attRepo,
			}
			tt.behaviour(tt.fields, tt.args)
			if err := logic.SetPayrollPeriod(tt.args.ctx, tt.args.start, tt.args.end); (err != nil) != tt.wantErr {
				t.Errorf("PayrollLogic.SetPayrollPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPayrollLogic_CalculatePay(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDeps := config.CommonDependencies{
		Config: config.InitConfig(context.Background()),
	}

	mockPayrollRepo := NewMockPayrollRepositoryInterface(ctrl)
	mockUserRepo := user.NewMockUserRepositoryInterface(ctrl)
	mockAttRepo := attendance.NewMockAttendanceRepositoryInterface(ctrl)
	type fields struct {
		deps        *config.CommonDependencies
		payrollRepo PayrollRepositoryInterface
		userRepo    user.UserRepositoryInterface
		attRepo     attendance.AttendanceRepositoryInterface
	}
	type args struct {
		ctx  context.Context
		data PayrollCalculationData
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      float64 // take home pay amount
		behaviour func(f fields, a args)
	}{
		// TODO: Add test cases.
		{
			name: "success calculate take home pay",
			fields: fields{
				deps:        &mockDeps,
				payrollRepo: mockPayrollRepo,
				userRepo:    mockUserRepo,
				attRepo:     mockAttRepo,
			},
			args: args{
				ctx: context.Background(),
				data: PayrollCalculationData{
					TotalWorkDay:       20,
					AttendanceCount:    20,
					OvertimeHoursCount: 8,
					Salary:             10000000,
					// overtime pay should be 62.500/hour
					Reimbursements: []Reimbursement{
						{
							Amount: 500000,
							Desc:   "buat jajan",
						},
						{
							Amount: 25000,
							Desc:   "buat ojol",
						},
					},
				},
			},
			want:      11025000,
			behaviour: func(f fields, a args) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := &PayrollLogic{
				deps:        tt.fields.deps,
				payrollRepo: tt.fields.payrollRepo,
				userRepo:    tt.fields.userRepo,
				attRepo:     tt.fields.attRepo,
			}
			if got := logic.CalculatePay(tt.args.ctx, tt.args.data); !reflect.DeepEqual(got.TakeHomePay, tt.want) {
				t.Errorf("PayrollLogic.CalculatePay() = %v, want %v", got.TakeHomePay, tt.want)
			}
		})
	}
}
