package main

import (
	pb "github.com/xdu31/test-server/pkg/pb"
	"github.com/xdu31/test-server/pkg/svc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	"github.com/infobloxopen/atlas-app-toolkit/requestid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// NewGRPCServer creates a new GRPC Server
func NewGRPCServer(logger *logrus.Logger, db *gorm.DB) (*grpc.Server, error) {
	interceptors := []grpc.UnaryServerInterceptor{
		grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger)),
		//logging.LogLevelInterceptor(logger.Level), // Request-scoped logging middleware
		// Request-Id interceptor
		requestid.UnaryServerInterceptor(),
		// collection operators middleware
		gateway.UnaryServerInterceptor(),
	}
	options := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors...))

	//interceptors_s := []grpc.StreamServerInterceptor{
	//}
	//options_s := []grpc.ServerOption{
	//	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(interceptors_s...))
	//}

	// create new gRPC grpcServer with middleware chain
	grpcServer := grpc.NewServer(options)

	// register service implementation with the grpcServer
	s, err := svc.NewBasicServer(db)
	if err != nil {
		return nil, err
	}
	pb.RegisterTestServerServer(grpcServer, s)

	// register ip service implementation with the grpcServer
	ips, err := svc.NewIpsServer(db)
	if err != nil {
		return nil, err
	}
	pb.RegisterIpsServer(grpcServer, ips)

	return grpcServer, nil
}
