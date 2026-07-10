package config

type LimiterItem struct {
	Rate       int
	Period     int
	LockPeriod int `json:",optional"`
}
