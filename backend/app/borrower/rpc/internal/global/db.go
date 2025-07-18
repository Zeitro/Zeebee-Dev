package global

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetGormDB() *gorm.DB {
	return db
}

func initDB(dsn string, isDebug bool) (err error) {
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		return err
	}

	if isDebug {
		gormDb = gormDb.Debug()
	}

	sqlDb, err := gormDb.DB()
	if err != nil {
		return err
	}
	err = sqlDb.Ping()
	if err != nil {
		return err
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(30)

	db = gormDb

	return nil
}
