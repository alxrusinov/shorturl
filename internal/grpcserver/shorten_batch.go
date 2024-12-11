package grpcserver

import (
	"context"
	"fmt"

	"github.com/alxrusinov/shorturl/internal/model"
	pb "github.com/alxrusinov/shorturl/pkg/grpchandler/proto"
)

func (g *GRPCServer) ShortenBatch(ctx context.Context, in *pb.ShortenBatchRequest) (*pb.ShortenBatchResponse, error) {
	var content []*model.StoreRecord

	for _, val := range in.Objects {

		shortenURL, err := g.Generator.GenerateRandomString()

		if err != nil {
			return nil, err
		}

		rec := &model.StoreRecord{
			UUID:          in.UserId,
			OriginalLink:  val.OriginalUrl,
			CorrelationID: val.CorrelationId,
			ShortLink:     shortenURL,
		}

		content = append(content, rec)
	}

	res, err := g.store.SetBatchLink(content)

	if err != nil {
		return nil, err
	}

	result := new(pb.ShortenBatchResponse)

	for _, val := range res {
		rec := &pb.ShortenItemResponse{
			CorrelationId: val.CorrelationID,
			ShortUrl:      fmt.Sprintf("%s/%s", g.responseAddr, val.ShortLink),
		}

		result.Objects = append(result.Objects, rec)
	}

	return result, nil
}
