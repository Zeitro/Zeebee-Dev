package model

type UsState struct {
	Base
	Name string `gorm:"column:name;varchar(50);not null;default:''" json:"name"`
	Code string `gorm:"column:code;char(2);not null;default:''" json:"code"`

	Cities []*UsCity `json:"cities" gorm:"foreignKey:StateId;references:Id"`
}

func (m UsState) TableName() string {
	return "us_states"
}
