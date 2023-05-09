package container

import (
	"boilerplate/internal/controler"
	"context"
	"errors"
	"testing"

	"go.uber.org/zap"
)

func TestApplication_GracefulShutdown(t *testing.T) {
	shutdownError := errors.New("error")
	ok := shutdownFunc(func(context.Context) error { return nil })
	err := shutdownFunc(func(context.Context) error { return shutdownError })

	type fields struct {
		Controlers controler.Controlers
		Log        *zap.Logger
		shutdowns  []shutdownFunc
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Should do a graceful shutdown",
			fields: fields{
				Controlers: controler.Controlers{},
				Log:        zap.NewNop(),
				shutdowns: []shutdownFunc{
					ok,
				},
			},
			args: args{ctx: context.Background()},
		},
		{
			name: "Should do a graceful shutdown with multiple shutdowns",
			fields: fields{
				Controlers: controler.Controlers{},
				Log:        zap.NewNop(),
				shutdowns: []shutdownFunc{
					ok,
					ok,
					ok,
				},
			},
			args: args{ctx: context.Background()},
		},
		{
			name: "Should return an error on shutdown",
			fields: fields{
				Controlers: controler.Controlers{},
				Log:        zap.NewNop(),
				shutdowns: []shutdownFunc{
					err,
				},
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
		{
			name: "Should return multiple errors on shutdown",
			fields: fields{
				Controlers: controler.Controlers{},
				Log:        zap.NewNop(),
				shutdowns: []shutdownFunc{
					err,
					err,
				},
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				Controlers: tt.fields.Controlers,
				Log:        tt.fields.Log,
				shutdowns:  tt.fields.shutdowns,
			}

			err := app.GracefulShutdown(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Application.GracefulShutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if errors.Is(err, shutdownError) {
					return
				}
				t.Errorf("Application.GracefulShutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
