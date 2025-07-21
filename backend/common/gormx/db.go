package gormx

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/service"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{
			Logger: NewGormLogger(logLevel),
		})
	default:
		return nil, errors.New("invalid db type: " + string(typ))
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
	err = sqlDb.Ping()
	if err != nil {
		return nil, err
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(30)

	return db, nil
}
