package grpcserver

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/alxrusinov/shorturl/internal/generator"
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	pb "github.com/alxrusinov/shorturl/pkg/grpchandler/proto"
	"github.com/stretchr/testify/mock"
)

func TestGRPCServer_GetOriginaltLink(t *testing.T) {
	testStore := mockstore.NewMockStore()
	testGenerator := generator.NewGenerator()

	type args struct {
		ctx context.Context
		in  *pb.GetOriginalRequest
	}
	tests := []struct {
		name    string
		g       *GRPCServer
		args    args
		want    *pb.GetOriginalResponse
		wantErr bool
	}{
		{
			name: "1# success",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in:  &pb.GetOriginalRequest{ShortenLink: "http://short"},
			},
			want: &pb.GetOriginalResponse{
				Link: "http://foo.bar",
			},
			wantErr: false,
		},
		{
			name: "2# error",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in:  &pb.GetOriginalRequest{ShortenLink: "http://short"},
			},
			want:    nil,
			wantErr: true,
		},
	}

	testStore.On("GetLink", mock.Anything).Return(&model.StoreRecord{OriginalLink: "http://foo.bar"}, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStore.On("GetLink", mock.Anything).Unset()
			if tt.wantErr {
				testStore.On("GetLink", mock.Anything).Return(new(model.StoreRecord), errors.New("error"))
			} else {
				testStore.On("GetLink", mock.Anything).Return(&model.StoreRecord{OriginalLink: "http://foo.bar"}, nil)
			}

			got, err := tt.g.GetOriginaltLink(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.GetOriginaltLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.GetOriginaltLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
