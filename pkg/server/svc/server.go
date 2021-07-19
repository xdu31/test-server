package svc

import (
	pb "github.com/xdu31/test-server/pkg/server/pb"
)

func NewIpsServer() pb.IpServiceServer {
	return &ipServer{}
}

type ipServer struct{}
