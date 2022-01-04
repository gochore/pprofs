package pprofs

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestEnableCapture(t *testing.T) {
	type args struct {
		options []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "invalid option",
			args: args{
				options: []Option{WithProfiles(nil)},
			},
			wantErr: ErrInvalidOption,
		},
		{
			name: "enabled",
			args: args{
				options: []Option{
					WithProfiles(
						CpuProfile().WithDuration(time.Second),
						HeapProfile(),
						MutexProfile(),
						BlockProfile().WithRate(1),
						GoroutineProfile(),
						ThreadcreateProfile(),
					),
					WithTrigger(NewIntervalTrigger(time.Second)),
					WithStorage(NewFileStorage("custom", "/tmp/pprofs/", time.Hour)),
					WithLogger(log.Default()),
				},
			},
			wantErr: nil,
		},
		{
			name: "already enabled",
			args: args{
				options: nil,
			},
			wantErr: ErrAlreadyEnabled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnableCapture(tt.args.options...)
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("EnableCapture() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	time.Sleep(5 * time.Second)
}
