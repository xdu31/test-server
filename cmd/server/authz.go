package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/dgrijalva/jwt-go"
	jwtgo "github.com/dgrijalva/jwt-go"

	tkgateway "github.com/infobloxopen/atlas-app-toolkit/gateway"

	"github.com/xdu31/test-server/pkg/util"
)

const (
	authKey = "authorization"
)

func AuthZUnaryServerInterceptor(authzAddress, appID string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		authStr, ok := tkgateway.Header(ctx, "Authorization")
		if !ok {
			fmt.Printf("in AuthZUnaryServerInterceptor not ok\n")
			return nil, util.PermissionDeniedErr("Authorization header is missing")
		}

		parts := strings.Split(authStr, "=")
		if len(parts) != 2 || parts[0] != "AccountID" {
			fields := strings.Fields(authStr)
			if len(fields) != 2 {
				fmt.Printf("in AuthZUnaryServerInterceptor not 2 fields\n")
				return nil, util.PermissionDeniedErr("Invalid Authorization header")
			}

			service, err := getJWTField(fields[1], "service")
			if err != nil {
				return nil, err
			}

			if service == "all" {
				return handler(ctx, req)
			}

			return nil, util.PermissionDeniedErr("Invalid Authorization header")
		}

		aid, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Printf("in AuthZUnaryServerInterceptor not 2 fields\n")
			return nil, util.PermissionDeniedErr("Invalid Account ID: %v", err)
		}

		ctx = context.WithValue(ctx, util.AccountIDKey, aid)
		ctx = updateMetaData(ctx, authKey, []string{"Bearer " + createJWTByAccount(aid)})

		return handler(ctx, req)
	}
}

func getJWTField(tokenStr string, tokenField string) (string, error) {
	parser := jwtgo.Parser{}
	token, _, err := parser.ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		return "", errors.New("unable to assert token as jwt.MapClaims")
	}
	jwtField, ok := claims[tokenField]
	if !ok {
		return "", errors.New("unable to get field from token")
	}
	return fmt.Sprint(jwtField), nil
}

func createJWTByAccount(aid int) string {
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, jwtgo.MapClaims{
		"account_id": aid,
		"username":   fmt.Sprintf("user_%d", aid),
		"exp":        time.Now().Add(time.Minute * 5).Unix(),
	})

	signedToken, _ := token.SignedString([]byte("secret"))
	return signedToken
}

func updateMetaData(ctx context.Context, key string, values []string) context.Context {
	var (
		md, _ = metadata.FromIncomingContext(ctx)
		updMD = metadata.New(map[string]string{}) // to be modified as demands by doc
	)

	for k, v := range md {
		updMD[k] = v
	}
	updMD[key] = values

	return metadata.NewIncomingContext(ctx, updMD)
}
