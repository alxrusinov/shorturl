package customerrors

import (
	"errors"
	"testing"
)

func TestDuplicateValueError_Unwrap(t *testing.T) {
	tests := []struct {
		name    string
		err     *DuplicateValueError
		wantErr bool
	}{
		{
			name: "1# success",
			err: &DuplicateValueError{
				Err: errors.New("error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.err.Unwrap(); (err != nil) != tt.wantErr {
				t.Errorf("DuplicateValueError.Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDuplicateValueError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *DuplicateValueError
		want string
	}{
		{
			name: "1# success",
			err: &DuplicateValueError{
				Err: errors.New("error"),
			},
			want: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("DuplicateValueError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
