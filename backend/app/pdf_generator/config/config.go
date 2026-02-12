package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	ApiConf ApiConf
}

type ApiConf struct {
	rest.RestConf
}
