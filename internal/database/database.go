package database

import (
	"errors"
	"net"
	"os"

	"github.com/stefanoschrs/proxymeister/pkg/types"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func Init() (db DB, err error) {
	sqlitePath := os.Getenv("DB_PATH")
	if sqlitePath == "" {
		sqlitePath = "file::memory:?cache=shared"
	}

	rawDB, err := gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
	if err != nil {
		return
	}

	db.DB = rawDB

	return
}

func (db DB) GetProxies() (proxies []types.Proxy, err error) {
	res := db.
		Order("updated_at DESC").
		Find(&proxies)
	if res.Error != nil {
		err = res.Error
		return
	}
	if len(proxies) == 0 {
		proxies = []types.Proxy{}
	}

	return
}

func (db DB) CreateProxy(p types.Proxy) (proxy types.Proxy, created bool, err error) {
	if net.ParseIP(p.Ip) == nil {
		err = errors.New(types.ErrInvalidIp)
		return
	}
	if p.Port <= 0 || p.Port > 65535 {
		err = errors.New(types.ErrInvalidPort)
		return
	}

	res := db.
		Where("ip = ? AND port = ?", p.Ip, p.Port).
		First(&proxy)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		err = res.Error
		return
	}
	if res.Error == nil {
		return
	}

	created = true
	res = db.Create(&p)
	if res.Error != nil {
		err = res.Error
		return
	}

	proxy = p
	return
}

func (db DB) DeleteProxy(id uint) (err error) {
	res := db.
		Where("id = ?", id).
		Delete(&types.Proxy{})
	if res.Error != nil {
		err = res.Error
		return
	}

	return
}

//func (db DB) UpdateAthleteSegmentStats(athleteId, segmentId, pr, efforts uint) (err error) {
//	entry := types.Entry{
//		SegmentId: segmentId,
//		AthleteId: athleteId,
//		PR:        pr,
//		Efforts:   efforts,
//	}
//	res := db.
//		Where(entry).
//		Assign(types.Entry{
//			PR:      pr,
//			Efforts: efforts,
//		}).
//		FirstOrCreate(&entry)
//	if res.Error != nil {
//		err = res.Error
//		return
//	}
//
//	return
//}
