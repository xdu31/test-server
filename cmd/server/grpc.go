package main

import (
	"context"

	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/infobloxopen/atlas-app-toolkit/errors"
	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	grpc_preprocess "github.com/infobloxopen/protoc-gen-preprocess/middleware"

	pb "github.com/xdu31/test-server/pkg/server/pb"
	"github.com/xdu31/test-server/pkg/server/svc"
)

// NewGRPCServer creates a new GRPC Server
func NewGRPCServer(logger *logrus.Logger) (*grpc.Server, error) {
	interceptors := []grpc.UnaryServerInterceptor{
		grpc_preprocess.UnaryServerInterceptor(),                    // preprocessing middleware
		grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger)), // logging middleware
		gateway.UnaryServerInterceptor(),                            // collection operators middleware
		AuthZUnaryServerInterceptor(
			viper.GetString("authz.host"),
			viper.GetString("authz.appid")),
	}

	options := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(interceptors...))
	// create new gRPC grpcServer with middleware chain
	grpcServer := grpc.NewServer(options)

	// register service implementation with the grpcServer
	pb.RegisterIpServiceServer(grpcServer, svc.NewIpsServer())

	return grpcServer, nil
}

func ValidateUnaryServerInterceptor(mapFuncs ...errors.MapFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res interface{}, err error) {
		mapper := errors.InitContainer().AddMapping(mapFuncs...)

		if validator, ok := req.(interface{ Validate() error }); ok {
			if err := validator.Validate(); err != nil {
				return nil, mapper.Map(ctx, err)
			}
		}

		return handler(ctx, req)
	}
}
