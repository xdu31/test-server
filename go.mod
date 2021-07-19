module github.com/xdu31/test-server

go 1.16

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190620084959-7cf5895f2711
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/go-openapi/spec v0.20.3
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/infobloxopen/atlas-app-toolkit v0.22.0
	github.com/infobloxopen/protoc-gen-preprocess v0.3.3
	github.com/jinzhu/gorm v1.9.14
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	google.golang.org/genproto v0.0.0-20210524171403-669157292da3
	google.golang.org/grpc v1.38.0
)

replace (
	github.com/infobloxopen/atlas-app-toolkit => github.com/infobloxopen/atlas-app-toolkit v0.18.2
	github.com/lyft/protoc-gen-validate => github.com/envoyproxy/protoc-gen-validate v0.0.7
)
