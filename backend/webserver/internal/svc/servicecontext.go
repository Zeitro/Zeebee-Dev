package svc

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"zeebee-backend/webserver/internal/config"
	"zeebee-backend/webserver/internal/model"
)

type ServiceContext struct {
	Config    config.Config
	UserModel *model.UserModel
	LoanModel *model.LoanModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := getGormDB(c.DataSource)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(db),
		LoanModel: model.NewLoanModel(db),
	}
}

func getGormDB(dsn string) (db *gorm.DB, err error) {
	return nil, err
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil

}
