package main

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func initFlags() {
	// define flag overrides
	pflag.String("grpc.address", "0.0.0.0", "adress of gRPC server")
	pflag.String("grpc.port", "9090", "port of gRPC server")
	pflag.Int("grpc.size", 1024*1024*4, "gRPC message size")

	pflag.String("gateway.address", "0.0.0.0", "address of gateway server")
	pflag.String("gateway.port", "8081", "port of gateway server")
	pflag.String("gateway.endpoint", "/api/test-server/v1/", "endpoint of gateway server")
	pflag.String("gateway.swaggerFile", "pkg/pb/service.swagger.json", "directory of swagger.json file")

	pflag.Bool("database.log", false, "database log")
	// DSN example: "postgres://postgres:postgres@postgres:5432/atlas_db?sslmode=disable"
	pflag.String("database.dsn", "", "DSN of the database")
	pflag.String("database.type", "postgres", "type of the database")
	pflag.String("database.address", "127.0.0.1", "address of the database")
	pflag.String("database.port", "5432", "port of the database")
	pflag.String("database.name", "test", "name of the database")
	pflag.String("database.user", "test_user", "database username")
	pflag.String("database.password", "test_pass", "database password")
	pflag.String("database.ssl", "disable", "database ssl mode")
	pflag.Int("database.retries", 3, "number of retries to connect the database")

	pflag.String("config.source", "deploy/", "directory of the configuration file")
	pflag.String("config.file", "", "directory of the configuration file")
	pflag.String("app.id", "testserver", "identifier for the application")

	pflag.String("log.level", "debug", "log level of application")

	pflag.Bool("authz.enable", false, "enable AuthZ service integration")
	pflag.String("authz.appid", "testserver", "enable AuthZ application ID")

	pflag.Bool("auditlog.enable", true, "enable logging of gRPC requests on Atlas audit log service")
	pflag.String("auditlog.address", "audit-logging.auditlog.svc:9090", "address and port of Atlas audit log service")
	pflag.Duration("audit.client.timeout", 10*time.Second, "timeout for audit log client")
	pflag.Int("audit.client.retry", 5, "retries to connect to audit log server")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(viper.GetString("config.source"))
	if viper.GetString("config.file") != "" {
		log.Printf("Serving from configuration file: %s", viper.GetString("config.file"))
		viper.SetConfigName(viper.GetString("config.file"))
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Cannot load configuration: %v", err)
		}
	} else {
		log.Printf("Serving from default values, environment variables, and/or flags")
	}
}
