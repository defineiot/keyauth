package mysql_test

import (
	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/domain"
	"github.com/defineiot/keyauth/dao/domain/mysql"
	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/internal/conf/mock"
)

func newTestStore() domain.Store {
	conf := mock.NewConfig()
	db, err := conf.GetDBConn()
	if err != nil {
		panic(err)
	}

	opt := &dao.Options{DB: db}
	store, err := mysql.NewDomainStore(opt)
	if err != nil {
		panic(err)
	}

	return store
}

type domainSuit struct {
	store domain.Store
	d1    *models.Domain
	d2    *models.Domain
}

func (d *domainSuit) TearDown() {
	if d.d1.ID != "" {
		d.store.DeleteDomainByID(d.d1.ID)
	}

	if d.d2.ID != "" {
		d.store.DeleteDomainByID(d.d2.ID)
	}

	d.store.Close()
}

func (d *domainSuit) SetUp() {
	d.d1 = &models.Domain{
		Name:           "unit_test_domain_name1",
		DisplayName:    "unit_test_domain_display_name1",
		LogoPath:       "/logo/to/path/unit_test_domain_logo1.png",
		Enabled:        true,
		Type:           models.Personal,
		Size:           "2000人以上",
		Location:       "中国,四川,成都",
		Address:        "环球中心 10F 1034",
		Industry:       "互联网",
		Fax:            "(+86 10)5992 0000",
		Phone:          "(+86 10)5992 8888",
		ContactsName:   "钟大俊",
		ContactsTitle:  "财务专员",
		ContactsMobile: "18188337463",
		ContactsEmail:  "18188337463@163.com",
		Owner:          "test_owner_01",
	}

	d.d2 = &models.Domain{
		Name:           "unit_test_domain_name2",
		DisplayName:    "unit_test_domain_display_name2",
		LogoPath:       "/logo/to/path/unit_test_domain_logo2.png",
		Enabled:        true,
		Type:           models.Personal,
		Size:           "2000人以上",
		Location:       "中国,四川,成都",
		Address:        "环球中心 10F 1034",
		Industry:       "互联网",
		Fax:            "(+86 10)5992 0000",
		Phone:          "(+86 10)5992 8888",
		ContactsName:   "钟大俊",
		ContactsTitle:  "财务专员",
		ContactsMobile: "18188337463",
		ContactsEmail:  "18188337463@163.com",
		Owner:          "test_owner_02",
	}

	d.store = newTestStore()

}
