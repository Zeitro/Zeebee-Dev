package global

import (
	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
	"zeebee-backend/common/gormx"
)

var db *gorm.DB

func GetGormDB() *gorm.DB {
	return db
}

func initDB(dsn string, isDebug bool) (err error) {
	gormDB, err := gormx.GetDB(dsn, isDebug, gormx.DBTypePgsql)
	if err != nil {
		return errorx.Wrap(err, "Init DB error")
	}
	db = gormDB
	return nil
}
