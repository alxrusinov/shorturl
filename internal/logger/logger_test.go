package logger

import (
	"os"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name string
		want zerolog.Logger
	}{
		{
			name: "1# success",
			want: zerolog.New(os.Stdout),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}
