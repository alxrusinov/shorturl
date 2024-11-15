package handler

import "testing"

func Test_createShortLink(t *testing.T) {

	tests := []struct {
		name    string
		host    string
		shorten string
		want    string
	}{
		{
			name:    "1# success",
			host:    "http://example.com",
			shorten: "123",
			want:    "http://example.com/123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createShortLink(tt.host, tt.shorten); got != tt.want {
				t.Errorf("createShortLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
