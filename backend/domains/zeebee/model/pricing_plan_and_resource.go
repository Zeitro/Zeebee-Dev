package model

type PricingPlanAndResource struct {
	Id            int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	PricingPlanId int64 `gorm:"column:pricing_plan_id;type:int4;not null;default:0" json:"pricing_plan_id"`
	ResourceId    int64 `gorm:"column:resource_id;type:int4;not null;default:0" json:"resource_id"`
}

func (PricingPlanAndResource) TableName() string {
	return "pricing_plan_and_resources"
}
