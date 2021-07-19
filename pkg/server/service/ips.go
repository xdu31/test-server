package service

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jinzhu/gorm"
	data "github.com/xdu31/test-server/pkg/server/db"
	"github.com/xdu31/test-server/pkg/server/pb"
	"github.com/xdu31/test-server/pkg/storage"
	"github.com/xdu31/test-server/pkg/util"
)

type IpService interface {
	Create(params map[string]interface{}, ip *data.Ip) (*data.Ip, error)
	Update(ip *data.Ip) (*data.Ip, error)
	Delete(ids []int) ([]*data.Ip, error)
	List(params map[string]interface{}) ([]*data.Ip, error)
	Get(id int) (*data.Ip, error)
}

type ipService struct {
	ai uint32
}

func NewIpService(ai uint32) IpService {
	return &ipService{ai: ai}
}

/*
 * @Description : List function returns all Ip objects for an
 * account given in AuthInfo.
 */
func (s *ipService) List(params map[string]interface{}) ([]*data.Ip, error) {
	var ips []*data.Ip
	params = util.EnsureParams(params)

	if s.ai != 0 {
		params["account_id"] = s.ai
	}

	if err := storage.WithTx(func(tx *gorm.DB) (err error) {
		ips, err = data.NewIpData(tx).Find(params)
		return err
	}); err != nil {
		return nil, err
	}

	return ips, nil
}

/*
 * @Description : Create function creates a new Ip object
 * for an account given in AuthInfo.
 */
func (s *ipService) Create(params map[string]interface{}, ip *data.Ip) (*data.Ip, error) {
	params = util.EnsureParams(params)
	params["account_id"] = s.ai

	ip.AccountId = int(s.ai)
	ip.Id = 0

	if err := storage.WithTx(func(tx *gorm.DB) (err error) {
		ip, err = data.NewIpData(tx).Create(params, ip)
		return err
	}); err != nil {
		return nil, err
	}

	return ip, nil
}

/*
 * @Description : Get function returns a Ip object by identifier
 * for an account given in AuthInfo.
 */
func (s *ipService) Get(id int) (*data.Ip, error) {
	var ip *data.Ip

	if err := storage.WithTx(func(tx *gorm.DB) (err error) {
		ip, err = data.NewIpData(tx).Get(int(s.ai), id)
		if err != nil {
			if err.Error() == gorm.ErrRecordNotFound.Error() {
				return util.NotFoundErr("Ip does not exist")
			}
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ip, nil
}

/*
 * @Description : Update function updates parameters of an Ip object
 * by identifier for an account given in AuthInfo.
 */
func (s *ipService) Update(ip *data.Ip) (*data.Ip, error) {
	ip.AccountId = int(s.ai)

	if err := storage.WithTx(func(tx *gorm.DB) (err error) {
		ip, err = data.NewIpData(tx).Update(ip)
		return err
	}); err != nil {
		return nil, err
	}

	return ip, nil
}

/*
 * @Description : Delete function deletes Ip objects by a list
 * of identifiers for an account given in AuthInfo.
 */
func (s *ipService) Delete(ids []int) ([]*data.Ip, error) {
	var (
		err error
		res []*data.Ip
	)

	if err = storage.WithTx(func(db *gorm.DB) (err error) {
		res, err = data.NewIpData(db).Delete(int(s.ai), ids)
		return err
	}); err != nil {
		return nil, err
	}

	return res, nil
}

/*
 * @Description : IpDBToTransport function performs
 * conversion of DB object representation to an API Ip.
 */
func IpDBToTransport(d *data.Ip) *pb.Ip {
	return &pb.Ip{
		Id:          int32(d.Id),
		IpAddress:   d.IpAddress,
		CreatedTime: &timestamp.Timestamp{Seconds: d.CreatedAt.Unix()},
		UpdatedTime: &timestamp.Timestamp{Seconds: d.UpdatedAt.Unix()},
	}
}
