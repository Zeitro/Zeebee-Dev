package global

import (
	"github.com/pkg/errors"

	"zeebee-backend/app/borrower/config"
	"zeebee-backend/common/redisx"
)

var redisDB *redisx.Redisx

func GetRedis() *redisx.Redisx {
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
		return errors.Wrap(err, "init redis error")
	}
	redisDB = rds
	return err
}
