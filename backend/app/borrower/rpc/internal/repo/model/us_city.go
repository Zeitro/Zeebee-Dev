package model

type UsCity struct {
	Base
	StateId int64  `gorm:"column:state_id;int4;not null;default:0" json:"state_id"`
	Name    string `gorm:"column:name;varchar(50);not null;default:''" json:"name"`
}

func (m UsCity) TableName() string {
	return "us_cities"
}
