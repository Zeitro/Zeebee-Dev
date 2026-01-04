package gormx

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/service"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"zeebee-backend/common/errorx"
)

type DBType string

const (
	DBTypePgsql DBType = "postgres"
	DBTypeMysql DBType = "mysql"
)

func GetDB(dsn string, isDebug bool, typ DBType) (db *gorm.DB, err error) {
	logLevel := service.DevMode
	if !isDebug {
		logLevel = service.ProMode
	}

	switch typ {
	case DBTypePgsql:
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{
			Logger: NewGormLogger(logLevel),
		})
	case DBTypeMysql:
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		}), &gorm.Config{
			Logger: NewGormLogger(logLevel),
		})
	default:
		return nil, errorx.New("Invalid DB Type: " + string(typ))
	}
	if err != nil {
		return nil, err
	}

	if isDebug {
		db = db.Debug()
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	err = sqlDb.PingContext(ctx)
	if err != nil {
		return nil, errorx.Wrap(err, "Ping DB Error")
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(30)

	return db, nil
}
