package grpcserver

import (
	"context"
	"errors"
	"testing"

	"github.com/alxrusinov/shorturl/internal/generator/mockgenerator"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	pb "github.com/alxrusinov/shorturl/pkg/grpchandler/proto"
	"github.com/stretchr/testify/assert"
)

func TestGRPCServer_GetShortenLink(t *testing.T) {
	testStore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	type args struct {
		ctx context.Context
		in  *pb.GetShortLinkRequest
	}
	tests := []struct {
		name    string
		g       *GRPCServer
		args    args
		want    *pb.GetShortLinkResponse
		wantErr bool
	}{
		{
			name: "1# success",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in:  &pb.GetShortLinkRequest{Link: "http://foo.bar", UserId: "1"},
			},
			want: &pb.GetShortLinkResponse{
				ShortenLink: "http://short",
			},
			wantErr: false,
		},
		{
			name: "2# error",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in:  &pb.GetShortLinkRequest{Link: "http://foo.bar", UserId: "1"},
			},
			want:    nil,
			wantErr: true,
		},
	}

	testGenerator.On("GenerateRandomString").Return("123", nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testGenerator.On("GenerateRandomString").Unset()
			if tt.wantErr {
				testGenerator.On("GenerateRandomString").Return("", errors.New("error"))
				got, err := tt.g.GetShortenLink(tt.args.ctx, tt.args.in)
				assert.NotNil(t, err)
				assert.Nil(t, got)
			} else {
				testGenerator.On("GenerateRandomString").Return("123", nil)
				got, err := tt.g.GetShortenLink(tt.args.ctx, tt.args.in)
				assert.NotNil(t, got)
				assert.Nil(t, err)
			}
		})
	}
}
