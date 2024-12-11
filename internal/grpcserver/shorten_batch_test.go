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

func TestGRPCServer_ShortenBatch(t *testing.T) {
	testStore := mockstore.NewMockStore()
	testGenerator := mockgenerator.NewMockGenerator()
	type args struct {
		ctx context.Context
		in  *pb.ShortenBatchRequest
	}
	tests := []struct {
		name    string
		g       *GRPCServer
		args    args
		want    *pb.ShortenBatchResponse
		wantErr bool
	}{
		{
			name: "1# success",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in: &pb.ShortenBatchRequest{UserId: "1", Objects: []*pb.ShortenItemRequest{{
					CorrelationId: "111",
					OriginalUrl:   "original",
				}}},
			},
			want: &pb.ShortenBatchResponse{
				Objects: []*pb.ShortenItemResponse{
					{
						ShortUrl:      ":7000/result",
						CorrelationId: "111",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "2# generator error",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in: &pb.ShortenBatchRequest{UserId: "1", Objects: []*pb.ShortenItemRequest{{
					CorrelationId: "111",
					OriginalUrl:   "original",
				}}},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "3# error",
			g:    NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16"),
			args: args{
				ctx: context.Background(),
				in: &pb.ShortenBatchRequest{UserId: "1", Objects: []*pb.ShortenItemRequest{{
					CorrelationId: "111",
					OriginalUrl:   "original",
				}}},
			},
			want:    nil,
			wantErr: true,
		},
	}

	testGenerator.On("GenerateRandomString").Return("result", nil)
	testStore.On("SetBatchLink", mock.Anything).Return([]*model.StoreRecord{{
		CorrelationID: "111",
		ShortLink:     "result",
	}}, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == tests[1].name {
				testGenerator.On("GenerateRandomString").Unset()
				testGenerator.On("GenerateRandomString").Return("", errors.New("generator error"))
			}

			if tt.name == tests[2].name {
				testGenerator.On("GenerateRandomString").Unset()
				testStore.On("SetBatchLink", mock.Anything).Unset()
				testGenerator.On("GenerateRandomString").Return("result", nil)
				testStore.On("SetBatchLink", mock.Anything).Return([]*model.StoreRecord{}, errors.New("error"))
			}

			got, err := tt.g.ShortenBatch(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.ShortenBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.ShortenBatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
