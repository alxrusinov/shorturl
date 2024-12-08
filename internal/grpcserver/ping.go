package grpcserver

import (
	"context"

	pb "github.com/alxrusinov/shorturl/internal/grpchandler/proto"
)

func (g *GRPCServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	var response pb.PingResponse

	err := g.store.Ping()

	if err != nil {
		return nil, err
	}

	return &response, nil
}
