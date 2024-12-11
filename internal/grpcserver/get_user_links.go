package grpcserver

import (
	"context"
	"fmt"

	pb "github.com/alxrusinov/shorturl/pkg/grpchandler/proto"
)

func (g *GRPCServer) GetUserLinks(ctx context.Context, in *pb.UserLinksRequest) (*pb.UserLinksResponse, error) {
	links, err := g.store.GetLinks(in.UserId)

	if err != nil {
		return nil, err
	}

	result := pb.UserLinksResponse{}

	for _, link := range links {
		newLink := &pb.UserLinkResponse{
			ShortUrl:    fmt.Sprintf("%s/%s", g.responseAddr, link.ShortLink),
			OriginalUrl: link.OriginalLink,
		}
		result.Objects = append(result.Objects, newLink)
	}

	return &result, nil
}
