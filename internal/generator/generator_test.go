package generator

import (
	"reflect"
	"testing"
)

func TestGenerator_GenerateRandomString(t *testing.T) {
	gen := NewGenerator()

	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "success",
			want:    24,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gen.GenerateRandomString()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("GenerateRandomString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_GenerateUserID(t *testing.T) {
	gen := NewGenerator()

	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "success",
			want:    24,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gen.GenerateUserID()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("GenerateUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGenerator(t *testing.T) {
	gen := &Generator{}
	tests := []struct {
		name string
		want *Generator
	}{
		{
			name: "1# success",
			want: gen,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGenerator(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}
