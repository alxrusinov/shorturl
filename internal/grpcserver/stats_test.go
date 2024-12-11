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
)

func TestGRPCServer_Stats(t *testing.T) {
	testStore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	type args struct {
		ctx context.Context
		in  *pb.StatsRequest
	}
	tests := []struct {
		name    string
		g       *GRPCServer
		args    args
		want    *pb.StatsResponse
		wantErr bool
	}{
		{
			name: "1# success",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "176.14.64.0/18"),
			args: args{
				ctx: context.Background(),
				in:  &pb.StatsRequest{XRealIp: "176.14.86.83"},
			},
			want:    &pb.StatsResponse{Urls: 12, Users: 10},
			wantErr: false,
		},
		{
			name: "2# empty real ip",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "176.14.64.0/18"),
			args: args{
				ctx: context.Background(),
				in:  &pb.StatsRequest{XRealIp: ""},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "3# untrusted ip",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "176.14.64.0/18"),
			args: args{
				ctx: context.Background(),
				in:  &pb.StatsRequest{XRealIp: "192.14.86.83"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "4# stat error",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "176.14.64.0/18"),
			args: args{
				ctx: context.Background(),
				in:  &pb.StatsRequest{XRealIp: "176.14.86.83"},
			},
			want:    nil,
			wantErr: true,
		},
	}

	testStore.On("GetStat").Return(&model.StatResponse{URLS: 12, Users: 10}, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == tests[3].name {
				testStore.On("GetStat").Unset()
				testStore.On("GetStat").Return(&model.StatResponse{URLS: 12, Users: 10}, errors.New("error"))
			}

			got, err := tt.g.Stats(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.Stats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.Stats() = %v, want %v", got, tt.want)
			}
		})
	}
}
