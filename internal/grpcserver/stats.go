package grpcserver

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/alxrusinov/shorturl/internal/grpchandler/proto"
	"github.com/alxrusinov/shorturl/internal/netutils"
)

func (g *GRPCServer) Stats(ctx context.Context, in *pb.StatsRequest) (*pb.StatsResponse, error) {
	if in.XRealIp == "" {
		return nil, errors.New("forbidden")
	}

	trusted, err := netutils.CheckSubnet(g.trustedSubnet, in.XRealIp)

	if !trusted || err != nil {
		return nil, fmt.Errorf("forbidden: %v", err)
	}

	res, err := g.store.GetStat()

	if err != nil {
		return nil, err
	}

	result := &pb.StatsResponse{
		Urls:  int32(res.URLS),
		Users: int32(res.Users),
	}

	return result, nil
}
