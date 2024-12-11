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

func TestGRPCServer_GetUserLinks(t *testing.T) {
	testStore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	type args struct {
		ctx context.Context
		in  *pb.UserLinksRequest
	}
	tests := []struct {
		name    string
		g       *GRPCServer
		args    args
		want    *pb.UserLinksResponse
		wantErr bool
	}{
		{
			name: "1# success",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in:  &pb.UserLinksRequest{UserId: "1"},
			},
			want: &pb.UserLinksResponse{
				Objects: []*pb.UserLinkResponse{
					{
						ShortUrl:    ":7000/Short",
						OriginalUrl: "Original",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "2# error",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in:  &pb.UserLinksRequest{UserId: "1"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	testStore.On("GetLinks", mock.Anything).Return([]model.StoreRecord{{OriginalLink: "Original", ShortLink: "Short"}}, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testStore.On("GetLinks", mock.Anything).Unset()

			if tt.wantErr {
				testStore.On("GetLinks", mock.Anything).Return([]model.StoreRecord{}, errors.New("error"))
			} else {
				testStore.On("GetLinks", mock.Anything).Return([]model.StoreRecord{{OriginalLink: "Original", ShortLink: "Short"}}, nil)

			}

			got, err := tt.g.GetUserLinks(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.GetUserLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.GetUserLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
