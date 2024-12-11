package grpcserver

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/alxrusinov/shorturl/internal/generator"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	pb "github.com/alxrusinov/shorturl/pkg/grpchandler/proto"
)

func TestGRPCServer_Ping(t *testing.T) {
	testStore := mockstore.NewMockStore()
	testGenerator := generator.NewGenerator()

	type args struct {
		ctx context.Context
		req *pb.PingRequest
	}
	tests := []struct {
		name    string
		g       *GRPCServer
		args    args
		want    *pb.PingResponse
		wantErr bool
	}{
		{
			name: "1# success",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				req: new(pb.PingRequest),
			},
			want:    new(pb.PingResponse),
			wantErr: false,
		},
		{
			name: "2# error",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				req: new(pb.PingRequest),
			},
			want:    nil,
			wantErr: true,
		},
	}

	testStore.On("Ping").Return(nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testStore.On("Ping").Unset()
			if tt.wantErr {
				testStore.On("Ping").Return(errors.New("error"))
			} else {
				testStore.On("Ping").Return(nil)
			}

			got, err := tt.g.Ping(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.Ping() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.Ping() = %v, want %v", got, tt.want)
			}
		})
	}
}
