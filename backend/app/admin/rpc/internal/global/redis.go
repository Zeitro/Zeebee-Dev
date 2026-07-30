package global

import (
	"zeebee-backend/app/admin/config"
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
		Prefix: cnf.Prefix,
	})
	if err != nil {
		return errorx.Wrap(err, "init redis error")
	}
	redisDB = rds
	return err
}
