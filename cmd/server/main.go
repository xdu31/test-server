package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	"github.com/infobloxopen/atlas-app-toolkit/server"

	"github.com/xdu31/test-server/pkg/server/pb"
	"github.com/xdu31/test-server/pkg/storage"
)

var (
	logger *logrus.Logger
)

func main() {
	initFlags()
	initLogger()

	storage.Open()
	defer storage.Close()

	logger.Info("Test service starts")

	if err := ServeExternal(logger); err != nil {
		logger.Fatal(err)
	}
}

func initLogger() {
	logger = logrus.StandardLogger()
	level := viper.GetString("log.level")

	ll, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Errorf("Invalid log level: %s", level)
		ll = logrus.InfoLevel
	}

	logger.SetLevel(ll)
}

// ServeExternal builds and runs the server that listens on ServerAddress and GatewayAddress
func ServeExternal(logger *logrus.Logger) error {
	grpcL, err := net.Listen("tcp", fmt.Sprintf("%s:%s", viper.GetString("grpc.address"), viper.GetString("grpc.port")))
	if err != nil {
		logger.Fatalln(err)
	}

	grpcServer, err := NewGRPCServer(logger)
	if err != nil {
		logger.Fatalln(err)
	}

	s, err := server.NewServer(
		server.WithGrpcServer(grpcServer),
		server.WithGateway(
			gateway.WithGatewayOptions(
				runtime.WithMetadata(gateway.MetadataAnnotator),
				runtime.WithForwardResponseOption(forwardResponseOption),
				runtime.WithMetadata(gateway.NewPresenceAnnotator("PUT", "POST", "PATCH")),
				runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
			),
			gateway.WithDialOptions(
				grpc.WithInsecure(),
				grpc.WithMaxMsgSize(viper.GetInt("grpc.size")),
				grpc.WithUnaryInterceptor(
					grpc_middleware.ChainUnaryClient(
						gateway.ClientUnaryInterceptor,
						gateway.PresenceClientInterceptor(),
					),
				),
			),
			gateway.WithServerAddress(fmt.Sprintf("%s:%s", viper.GetString("grpc.address"), viper.GetString("grpc.port"))),
			gateway.WithEndpointRegistration(viper.GetString("gateway.endpoint"),
				pb.RegisterIpServiceHandlerFromEndpoint),
		),
		server.WithHandler("/swagger/", swaggerHandler(viper.GetString("gateway.swaggerFile"))),
	)
	if err != nil {
		logger.Fatalln(err)
	}

	httpL, err := net.Listen("tcp", fmt.Sprintf("%s:%s", viper.GetString("gateway.address"), viper.GetString("gateway.port")))
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Printf("serving gRPC at %s:%s", viper.GetString("grpc.address"), viper.GetString("grpc.port"))
	logger.Printf("serving http at %s:%s", viper.GetString("gateway.address"), viper.GetString("gateway.port"))

	return s.Serve(grpcL, httpL)
}

func forwardResponseOption(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	return nil
}
