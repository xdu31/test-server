package util

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SetOK(ctx context.Context, format string, args ...interface{}) error {
	return gateway.SetStatus(ctx, status.New(codes.OK, fmt.Sprintf(format, args...)))
}

func EnsureParams(params map[string]interface{}) map[string]interface{} {
	if params == nil {
		return map[string]interface{}{}
	}
	return params
}

/*
this function will compress the ipv6 addresses
Example:
ips := []string{"2001:db8:0:0:0:0:0:0","2001:db8"}
will result into below :
"[2001:db8,2001:db8]"
*/
func CompressIPV6(ips []sql.NullString) []sql.NullString {
	for i, item := range ips {
		if CheckIPv6Addr(item.String) {
			tmpItemI := net.ParseIP(item.String).To16().String()
			for j := i + 1; j < len(ips); j++ {
				tmpItemJ := net.ParseIP(ips[j].String).To16().String()
				if tmpItemI == tmpItemJ {
					ips[j].String = tmpItemI
					ips[i].String = tmpItemI
				}
			}
		}
	}
	return ips
}

//InStringSlice.. this will check whether the value is present in string slice or not.
func InStringSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
