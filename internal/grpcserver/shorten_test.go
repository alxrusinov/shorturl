package grpcserver

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/alxrusinov/shorturl/internal/generator/mockgenerator"
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	pb "github.com/alxrusinov/shorturl/pkg/grpchandler/proto"
	"github.com/stretchr/testify/mock"
)

func TestGRPCServer_Shorten(t *testing.T) {
	testStore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	type args struct {
		ctx context.Context
		in  *pb.ShortenRequest
	}
	tests := []struct {
		name    string
		g       *GRPCServer
		args    args
		want    *pb.ShortenResponse
		wantErr bool
	}{
		{
			name: "1# success",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in:  &pb.ShortenRequest{UserId: "1", Url: "original"},
			},
			want: &pb.ShortenResponse{
				Result: ":7000/result",
			},
			wantErr: false,
		},
		{
			name: "2# generator error",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in:  &pb.ShortenRequest{UserId: "1", Url: "original"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "3# error",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in:  &pb.ShortenRequest{UserId: "1", Url: "original"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	testGenerator.On("GenerateRandomString").Return("result", nil)
	testStore.On("SetLink", mock.Anything).Return(&model.StoreRecord{ShortLink: "result"}, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == tests[1].name {
				testGenerator.On("GenerateRandomString").Unset()
				testGenerator.On("GenerateRandomString").Return("r", errors.New("generator error"))
			}

			if tt.name == tests[2].name {
				testGenerator.On("GenerateRandomString").Unset()
				testStore.On("SetLink", mock.Anything).Unset()

				testGenerator.On("GenerateRandomString").Return("result", nil)
				testStore.On("SetLink", mock.Anything).Return(&model.StoreRecord{ShortLink: "result"}, errors.New("error"))
			}

			got, err := tt.g.Shorten(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.Shorten() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.Shorten() = %v, want %v", got, tt.want)
			}
		})
	}
}
