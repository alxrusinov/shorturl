package handler

import (
	"testing"
)

func Test_checkContentType(t *testing.T) {
	type args struct {
		values []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1# true",
			args: args{
				values: []string{"text/html", "application/json"},
			},
			want: true,
		},
		{
			name: "1# false",
			args: args{
				values: []string{"text", "json"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkContentType(tt.args.values); got != tt.want {
				t.Errorf("checkContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkGzip(t *testing.T) {
	type args struct {
		values []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1# true",
			args: args{
				values: []string{"gzip", "json"},
			},
			want: true,
		},
		{
			name: "1# false",
			args: args{
				values: []string{"text", "json"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkGzip(tt.args.values); got != tt.want {
				t.Errorf("checkGzip() = %v, want %v", got, tt.want)
			}
		})
	}
}
