package grpcserver

import (
	"context"

	"github.com/alxrusinov/shorturl/internal/model"
	pb "github.com/alxrusinov/shorturl/pkg/grpchandler/proto"
)

// APIDeleteLinks - deletes links for user
func (g *GRPCServer) APIDeleteLinks(ctx context.Context, in *pb.DeleteLinkRequest) (*pb.DeleteResponse, error) {
	var batch []model.StoreRecord

	for _, val := range in.Links {
		batch = append(batch, model.StoreRecord{
			UUID:      in.UserId,
			ShortLink: val,
		})
	}

	g.deleteChan <- batch

	return nil, nil
}
