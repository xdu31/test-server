package svc

import (
	"context"

	"github.com/infobloxopen/atlas-app-toolkit/gateway"
	"github.com/xdu31/test-server/pkg/server/db"
	"github.com/xdu31/test-server/pkg/server/pb"
	"github.com/xdu31/test-server/pkg/server/service"
	"github.com/xdu31/test-server/pkg/util"
)

const (
	ipResourceType  string = "IP"
	ipCreateMessage string = "IP has been created"
	ipUpdateMessage string = "IP has been updated"
	ipDeleteMessage string = "IP has been deleted"

	opUpdate = "UPDATE"
)

func newIpService(ctx context.Context) (service.IpService, error) {
	ai, err := util.GetAccountID(ctx)
	if err != nil {
		return nil, err
	}
	return service.NewIpService(ai), nil
}

func IpTransportToDB(d *pb.Ip) *db.Ip {
	return &db.Ip{
		Id:        int(d.Id),
		IpAddress: d.IpAddress,
	}
}

/*
 * @Description : List function implements retrieval of IPs
 */
func (s *ipServer) List(ctx context.Context, ip *pb.ListIpsRequest) (*pb.ListIpsResponse, error) {
	srv, err := newIpService(ctx)
	if err != nil {
		return nil, util.ErrStatus(err)
	}

	ipList, err := srv.List(nil)
	if err != nil {
		return nil, util.ErrStatus(err)
	}

	data := []*pb.Ip{}
	for _, d := range ipList {
		data = append(data, service.IpDBToTransport(d))
	}

	return &pb.ListIpsResponse{Results: data}, util.SetOK(ctx, "Found %d IP(s)", len(data))
}

/*
 * @Description : Delete function implements deletion of IPs
 */
func (s *ipServer) Delete(ctx context.Context, ip *pb.DeleteIpRequest) (*pb.DeleteIpResponse, error) {
	srv, err := newIpService(ctx)
	if err != nil {
		return nil, util.ErrStatus(err)
	}

	if len(ip.Ids) == 0 {
		return nil, util.ErrStatus(util.FailedPreconditionErr("IP 'ids' can't be empty"))
	}

	idList := make([]int, len(ip.Ids))
	for i, id := range ip.Ids {
		idList[i] = int(id)
	}

	_, err = srv.Delete(idList)
	if err != nil {
		return nil, util.ErrStatus(err)
	}

	return &pb.DeleteIpResponse{}, gateway.SetDeleted(ctx, ipDeleteMessage)
}

/*
 * @Description : Update function implements update of IP
 */
func (s *ipServer) Update(ctx context.Context, in *pb.UpdateIpRequest) (*pb.UpdateIpResponse, error) {
	srv, err := newIpService(ctx)
	if err != nil {
		return nil, util.ErrStatus(err)
	}

	ip := in.Payload

	data, err := srv.Update(IpTransportToDB(ip))
	if err != nil {
		return nil, util.ErrStatus(err)
	}

	return &pb.UpdateIpResponse{Result: service.IpDBToTransport(data)},
		gateway.SetUpdated(ctx, ipUpdateMessage)
}

/*
 * @Description : Read function implements read of IP
 */
func (s *ipServer) Read(ctx context.Context, req *pb.ReadIpRequest) (*pb.ReadIpResponse, error) {
	srv, err := newIpService(ctx)
	if err != nil {
		return nil, util.ErrStatus(err)
	}

	d, err := srv.Get(int(req.Id))
	if err != nil {
		return nil, util.ErrStatus(err)
	}

	return &pb.ReadIpResponse{Result: service.IpDBToTransport(d)}, util.SetOK(ctx, "Found IP")
}

/*
 * @Description : Create function implements creation of IP
 */
func (s *ipServer) Create(ctx context.Context, in *pb.CreateIpRequest) (*pb.CreateIpResponse, error) {
	srv, err := newIpService(ctx)
	if err != nil {
		return nil, util.ErrStatus(err)
	}

	ip := in.Payload

	data, err := srv.Create(nil, IpTransportToDB(ip))
	if err != nil {
		return nil, util.ErrStatus(err)
	}

	return &pb.CreateIpResponse{Result: service.IpDBToTransport(data)},
		util.SetOK(ctx, ipCreateMessage)
}
