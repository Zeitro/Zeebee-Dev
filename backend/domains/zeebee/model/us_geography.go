package model

type UsGeography struct {
	Base
	StateName  string `gorm:"column:state_name;type:varchar(50);not null;default ''" json:"state_name"`
	StateUsps  string `gorm:"column:state_usps;type:varchar(10);not null;default ''" json:"state_usps"`
	CityName   string `gorm:"column:city_name;type:varchar(50);not null;default ''" json:"city_name"`
	CountyName string `gorm:"column:county_name;type:varchar(10);not null;default ''" json:"county_name"`
}

func (m *UsGeography) TableName() string {
	return "us_geographies"
}
