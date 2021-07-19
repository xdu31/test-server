package db

import (
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/xdu31/test-server/pkg/util"
)

type Ip struct {
	Id        int       `gorm:"column:id"`
	AccountId int       `gorm:"column:account_id"`
	IpAddress string    `gorm:"column:ip_address"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

const (
	IpTableName = "ips"
	SqlInsertIp = "INSERT INTO " + IpTableName + " (ip_address, account_id, created_at, updated_at) "
)

type IpData interface {
	Find(params map[string]interface{}) ([]*Ip, error)
	Create(params map[string]interface{}, ip *Ip) (*Ip, error)
	Get(accountId int, id int) (*Ip, error)
	Update(ip *Ip) (*Ip, error)
	Delete(accountId int, ids []int) ([]*Ip, error)
}

type ipData struct {
	db *gorm.DB
}

func NewIpData(db *gorm.DB) IpData {
	return &ipData{db: db}
}

/*
 * @Description : Find function returns list of Ip objects.
 */
func (d *ipData) Find(params map[string]interface{}) ([]*Ip, error) {
	params = util.EnsureParams(params)
	var ips []*Ip

	query := d.db.Select("id,account_id,ip_address,created_at,updated_at").Where(params).Order("id")
	if err := query.Find(&ips).Error; err != nil {
		return nil, err
	}

	return ips, nil
}

/*
 * @Description : Create function creates a new Ip object.
 */
func (d *ipData) Create(params map[string]interface{}, ip *Ip) (*Ip, error) {
	insertIntoIpsRow := d.db.Raw(SqlInsertIp, ip.IpAddress, ip.AccountId, ip.AccountId).Row()
	err := insertIntoIpsRow.Scan(&ip.Id, &ip.CreatedAt, &ip.UpdatedAt)

	return ip, err
}

/*
 * @Description : Get function returns Ip object by id.
 */
func (d *ipData) Get(accountId int, id int) (*Ip, error) {
	var ip Ip

	params := map[string]interface{}{
		"id":         id,
		"account_id": accountId,
	}

	if err := d.db.Select("id,account_id,ip_address,created_at,updated_at").Where(params).First(&ip).Error; err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, util.NotFoundErr("Ip does not exist")
		}
		return nil, err
	}

	return &ip, nil
}

/*
 * @Description : Update function updates a Ip object by identifier for a given account.
 */
func (d *ipData) Update(ip *Ip) (*Ip, error) {
	conditions := map[string]interface{}{
		"id":         ip.Id,
		"account_id": ip.AccountId,
	}

	params := map[string]interface{}{}

	resp := d.db.Model(&Ip{}).Where(conditions).Update(params)
	err := resp.Error
	if err != nil {
		return nil, err
	}

	if resp.RowsAffected != 1 {
		return nil, util.NotFoundErr("Ip does not exist")
	}

	ip, err = d.Get(ip.AccountId, ip.Id)
	if err != nil {
		return nil, err
	}

	return ip, nil
}

/*
 * @Description : Delete function deletes Ip objects by list of identifiers.
 */
func (d *ipData) Delete(accountId int, ids []int) ([]*Ip, error) {
	var res []*Ip

	err := d.db.Raw("DELETE FROM ips WHERE id IN (?) AND account_id=? RETURNING id, name", ids, accountId).Scan(&res).Error
	if err != nil {
		return nil, err
	}

	removed := map[int]struct{}{}
	for _, item := range res {
		removed[item.Id] = struct{}{}
	}

	if len(removed) != len(ids) {
		notExist := []string{}

		for _, id := range ids {
			if _, ok := removed[id]; !ok {
				notExist = append(notExist, strconv.Itoa(int(id)))
			}
		}

		return nil, util.FailedPreconditionErr("Non-existent: %s", strings.Join(notExist, ","))
	}

	return res, nil
}
