package grpcserver

import (
	"net"

	"github.com/alxrusinov/shorturl/internal/model"
	pb "github.com/alxrusinov/shorturl/pkg/grpchandler/proto"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

// Store - store for grpc
type Store interface {
	GetLink(arg *model.StoreRecord) (*model.StoreRecord, error)
	SetLink(arg *model.StoreRecord) (*model.StoreRecord, error)
	SetBatchLink(arg []*model.StoreRecord) ([]*model.StoreRecord, error)
	Ping() error
	GetLinks(userID string) ([]model.StoreRecord, error)
	DeleteLinks(shorts [][]model.StoreRecord) error
	GetStat() (*model.StatResponse, error)
}

// Type Generator is a type for generator
type Generator interface {
	GenerateRandomString() (string, error)
	GenerateUserID() (string, error)
}

// GRPCServer - typ of grpc server
type GRPCServer struct {
	pb.UnimplementedHandlerServer
	store         Store
	addr          string
	responseAddr  string
	trustedSubnet string
	deleteChan    chan []model.StoreRecord
	Generator     Generator
}

// Run - method of GRPCServer for runnning app
func Run(g *GRPCServer) error {
	server, err := net.Listen("tcp", g.addr)

	if err != nil {
		return err
	}
	s := grpc.NewServer()

	pb.RegisterHandlerServer(s, g)

	err = s.Serve(server)

	if err != nil {
		return err
	}

	return nil

}

func (g *GRPCServer) GetDelChan() <-chan []model.StoreRecord {
	return g.deleteChan
}

// NewGRPCServer creates GRPCServer
func NewGRPCServer(store Store, addr string, generator Generator, responseAddr string, trustedSubnet string) *GRPCServer {
	server := &GRPCServer{
		store:         store,
		addr:          addr,
		deleteChan:    make(chan []model.StoreRecord),
		Generator:     generator,
		responseAddr:  responseAddr,
		trustedSubnet: trustedSubnet,
	}

	return server
}
