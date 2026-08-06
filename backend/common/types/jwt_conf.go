package types

type JwtConf struct {
	AccessSecret  string `json:",default=accessSecret"`
	AccessExpire  int64  `json:",default=3600"`
	RefreshSecret string `json:",default=refreshSecret"`
	RefreshExpire int64  `json:",default=86400"`
}
