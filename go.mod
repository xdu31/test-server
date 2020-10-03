module github.com/xdu31/test-server

go 1.14

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-delve/delve v1.4.1 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.9.2
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645 // indirect
	github.com/infobloxopen/atlas-app-toolkit v0.22.0
	github.com/infobloxopen/protoc-gen-gorm v0.20.0
	github.com/jinzhu/gorm v1.9.14
	github.com/lyft/protoc-gen-validate v0.0.13
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5 // indirect
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1 // indirect
	google.golang.org/genproto v0.0.0-20191108220845-16a3f7862a1a
	google.golang.org/grpc v1.30.0
)

replace (
	github.com/infobloxopen/atlas-app-toolkit => github.com/infobloxopen/atlas-app-toolkit v0.18.2
	github.com/lyft/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.0.7
)
