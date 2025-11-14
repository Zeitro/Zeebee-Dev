package global

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"zeebee-backend/common/gormx"
)

var db *gorm.DB

func GetGormDB() *gorm.DB {
	return db
}

func initDB(dsn string, isDebug bool) (err error) {

	gormDB, err := gormx.GetDB(dsn, isDebug, gormx.DBTypePgsql)
	if err != nil {
		return errors.Wrap(err, "Init DB Error")
	}
	db = gormDB
	return nil
}
