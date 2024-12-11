package grpcserver

import (
	"context"

	pb "github.com/alxrusinov/shorturl/pkg/grpchandler/proto"
)

func (g *GRPCServer) GetShortenLink(ctx context.Context, in *pb.GetShortLinkRequest) (*pb.GetShortLinkResponse, error) {
	shortenURL, err := g.Generator.GenerateRandomString()

	if err != nil {
		return nil, err
	}

	result := &pb.GetShortLinkResponse{
		ShortenLink: shortenURL,
	}

	return result, nil
}
