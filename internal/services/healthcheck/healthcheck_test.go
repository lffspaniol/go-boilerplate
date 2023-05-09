package healthcheck_test

import (
	"boilerplate/internal/services/healthcheck"
	"boilerplate/internal/services/healthcheck/mock"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestAlive_Readiness(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		mock    func(mock *mock.MockDepencencie)
		args    args
		wantErr bool
	}{
		{
			name: "Should not return an error",
			mock: func(mock *mock.MockDepencencie) {
				mock.EXPECT().Healthcheck(gomock.Any()).Return(nil)
			},
			args: args{ctx: context.Background()},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			mock := mock.NewMockDepencencie(ctrl)
			tt.mock(mock)
			alive := healthcheck.New()
			if err := alive.Readiness(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Alive.Readiness() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
