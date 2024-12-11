package mockgenerator

import (
	"errors"
	"testing"
)

func TestMockGenerator_GenerateUserID(t *testing.T) {
	tests := []struct {
		name    string
		mg      *MockGenerator
		want    string
		wantErr bool
	}{
		{
			name:    "1# success",
			mg:      NewMockGenerator(),
			want:    "111",
			wantErr: false,
		},
		{
			name:    "1# error",
			mg:      NewMockGenerator(),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				tt.mg.On("GenerateUserID").Return("", errors.New("err"))
			} else {
				tt.mg.On("GenerateUserID").Return(tt.want, nil)
			}

			got, err := tt.mg.GenerateUserID()
			if (err != nil) != tt.wantErr {
				t.Errorf("MockGenerator.GenerateUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MockGenerator.GenerateUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}
