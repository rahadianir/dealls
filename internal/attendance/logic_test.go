package attendance

import (
	"context"
	"testing"
	"time"

	"github.com/rahadianir/dealls/internal/config"
	"go.uber.org/mock/gomock"
)

func TestAttendanceLogic_SubmitAttendance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockAttendanceRepositoryInterface(ctrl)
	mockDeps := config.CommonDependencies{
		Config: config.InitConfig(context.Background()),
	}
	tudei, err := time.Parse(time.RFC3339, "2025-06-11T06:29:44+07:00")
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		deps    *config.CommonDependencies
		attRepo AttendanceRepositoryInterface
		today   time.Time
	}
	type args struct {
		ctx       context.Context
		userID    string
		timestamp string
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
			name: "success submit attendance",
			fields: fields{
				deps:    &mockDeps,
				attRepo: mockRepo,
				today:   tudei,
			},
			args: args{
				ctx:       context.Background(),
				userID:    "user-id",
				timestamp: "2025-06-11T06:29:44+07:00",
			},
			wantErr: false,
			behaviour: func(f fields, a args) {
				mockRepo.EXPECT().SubmitAttendance(gomock.Any(), "user-id", gomock.Any()).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := &AttendanceLogic{
				deps:    tt.fields.deps,
				attRepo: tt.fields.attRepo,
				today:   tt.fields.today,
			}
			tt.behaviour(tt.fields, tt.args)
			if err := logic.SubmitAttendance(tt.args.ctx, tt.args.userID, tt.args.timestamp); (err != nil) != tt.wantErr {
				t.Errorf("AttendanceLogic.SubmitAttendance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAttendanceLogic_SubmitOvertime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockAttendanceRepositoryInterface(ctrl)
	mockDeps := config.CommonDependencies{
		Config: config.InitConfig(context.Background()),
	}
	tudei, err := time.Parse(time.RFC3339, "2025-06-11T06:29:44+07:00")
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		deps    *config.CommonDependencies
		attRepo AttendanceRepositoryInterface
		today   time.Time
	}
	type args struct {
		ctx                       context.Context
		userID                    string
		hourCount                 int
		finishedOvertimeTimestamp string
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
			name: "success submit overtime",
			fields: fields{
				deps:    &mockDeps,
				attRepo: mockRepo,
				today:   tudei,
			},
			args: args{
				ctx:                       context.Background(),
				userID:                    "user-id",
				hourCount:                 2,
				finishedOvertimeTimestamp: "2025-06-11T06:29:44+07:00",
			},
			wantErr: false,
			behaviour: func(f fields, a args) {
				mockRepo.EXPECT().GetUserOvertimeByTime(gomock.Any(), "user-id", gomock.Any()).Return(0, nil)
				mockRepo.EXPECT().SubmitOvertime(gomock.Any(), "user-id", 2, gomock.Any()).Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := &AttendanceLogic{
				deps:    tt.fields.deps,
				attRepo: tt.fields.attRepo,
				today:   tt.fields.today,
			}
			tt.behaviour(tt.fields, tt.args)
			if err := logic.SubmitOvertime(tt.args.ctx, tt.args.userID, tt.args.hourCount, tt.args.finishedOvertimeTimestamp); (err != nil) != tt.wantErr {
				t.Errorf("AttendanceLogic.SubmitOvertime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAttendanceLogic_SubmitReimbursement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockAttendanceRepositoryInterface(ctrl)
	mockDeps := config.CommonDependencies{
		Config: config.InitConfig(context.Background()),
	}
	tudei, err := time.Parse(time.RFC3339, "2025-06-11T06:29:44+07:00")
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		deps    *config.CommonDependencies
		attRepo AttendanceRepositoryInterface
		today   time.Time
	}
	type args struct {
		ctx    context.Context
		userID string
		amount float64
		desc   string
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
			name: "success submit reimbursement",
			fields: fields{
				deps:    &mockDeps,
				attRepo: mockRepo,
				today:   tudei,
			},
			args: args{
				ctx:    context.Background(),
				userID: "user-id",
				amount: 100,
				desc:   "desc",
			},
			wantErr: false,
			behaviour: func(f fields, a args) {
				mockRepo.EXPECT().SubmitReimbursement(gomock.Any(), "user-id", float64(100), "desc").Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := &AttendanceLogic{
				deps:    tt.fields.deps,
				attRepo: tt.fields.attRepo,
				today:   tt.fields.today,
			}
			tt.behaviour(tt.fields, tt.args)
			if err := logic.SubmitReimbursement(tt.args.ctx, tt.args.userID, tt.args.amount, tt.args.desc); (err != nil) != tt.wantErr {
				t.Errorf("AttendanceLogic.SubmitReimbursement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
