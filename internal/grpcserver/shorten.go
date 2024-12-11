package grpcserver

import (
	"context"
	"fmt"

	"github.com/alxrusinov/shorturl/internal/model"
	pb "github.com/alxrusinov/shorturl/pkg/grpchandler/proto"
)

func (g *GRPCServer) Shorten(ctx context.Context, in *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	shortenURL, err := g.Generator.GenerateRandomString()

	if err != nil {
		return nil, err
	}

	links := &model.StoreRecord{
		ShortLink:    shortenURL,
		OriginalLink: in.Url,
		UUID:         in.UserId,
	}

	res, err := g.store.SetLink(links)

	if err != nil {
		return nil, err
	}

	result := &pb.ShortenResponse{
		Result: fmt.Sprintf("%s/%s", g.responseAddr, res.ShortLink),
	}

	return result, nil
}
