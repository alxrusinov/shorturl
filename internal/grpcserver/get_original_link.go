package grpcserver

import (
	"context"

	pb "github.com/alxrusinov/shorturl/internal/grpchandler/proto"
	"github.com/alxrusinov/shorturl/internal/model"
)

func (g *GRPCServer) GetOriginaltLink(ctx context.Context, in *pb.GetOriginalRequest) (*pb.GetOriginalResponse, error) {

	link := &model.StoreRecord{
		ShortLink: in.ShortenLink,
	}

	res, err := g.store.GetLink(link)

	if err != nil {
		return nil, err
	}
	result := &pb.GetOriginalResponse{
		Link: res.OriginalLink,
	}

	return result, nil
}
