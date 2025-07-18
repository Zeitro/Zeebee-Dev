package types

type JwtAuth struct {
	AccessSecret string `json:",default=accessSecret"`
	AccessExpire int64  `json:",default=7200"`
}
