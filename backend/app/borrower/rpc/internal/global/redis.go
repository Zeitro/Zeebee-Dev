package global

import (
	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/errorx"
	"zeebee-backend/common/redisx"
)

var redisDB *redisx.Client

func GetRedis() *redisx.Client {
	return redisDB
}

func initRedis(cnf config.RedisConf) (err error) {
	rds, err := redisx.GetRedis(redisx.Opt{
		Host:   cnf.Host,
		Type:   cnf.Type,
		Pass:   cnf.Pass,
		Tls:    cnf.Tls,
		DB:     cnf.DB,
		Prefix: cnf.Prefix,
	})
	if err != nil {
		return errorx.Wrap(err, "Init Redis Error")
	}
	redisDB = rds
	return err
}
