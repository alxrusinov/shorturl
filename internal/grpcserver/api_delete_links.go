package grpcserver

import (
	"context"

	pb "github.com/alxrusinov/shorturl/internal/grpchandler/proto"
	"github.com/alxrusinov/shorturl/internal/model"
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

	g.DeleteChan <- batch

	return nil, nil
}
