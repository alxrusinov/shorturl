package grpcserver

import (
	"context"
	"reflect"
	"testing"

	"github.com/alxrusinov/shorturl/internal/generator"
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/alxrusinov/shorturl/internal/store/mockstore"
	pb "github.com/alxrusinov/shorturl/pkg/grpchandler/proto"
	"github.com/stretchr/testify/mock"
)

func TestGRPCServer_APIDeleteLinks(t *testing.T) {
	testStore := mockstore.NewMockStore()
	testGenerator := generator.NewGenerator()

	testStore.On("DeleteLinks", mock.Anything).Return(nil)

	gServer := NewGRPCServer(testStore, ":8000", testGenerator, ":7000", "196.168.0.0/16")

	go func() {
		var batch [][]model.StoreRecord

		delChan := gServer.GetDelChan()

		for val := range delChan {
			batch = append(batch, val)
			testStore.DeleteLinks(batch)

			batch = batch[0:0]
		}
	}()

	type args struct {
		ctx context.Context
		in  *pb.DeleteLinkRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *pb.DeleteResponse
		wantErr bool
	}{
		{
			name: "1# success",
			args: args{
				ctx: context.Background(),
				in: &pb.DeleteLinkRequest{
					UserId: "1",
					Links:  []string{"del1", "del1"},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gServer.APIDeleteLinks(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GRPCServer.APIDeleteLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GRPCServer.APIDeleteLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
