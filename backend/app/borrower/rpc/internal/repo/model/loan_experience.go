package model

import (
	"database/sql"

	"zeebee-backend/common/constants"
)

// CREATE TABLE "public"."loan_experiences" (
//  "id" int4 NOT NULL DEFAULT nextval('loan_experiences_id_seq'::regclass),
//  "created_at" timestamptz(6) NOT NULL DEFAULT now(),
//  "updated_at" timestamptz(6) NOT NULL DEFAULT now(),
//  "is_del" int2 NOT NULL DEFAULT 0,
//  "loan_id" int4 NOT NULL DEFAULT 0,
//  "city" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
//  "state" varchar(20) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
//  "ownership_on_title" varchar(100) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
//  "property_type" varchar(50) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
//  "purchase_amount" numeric(10,2) NOT NULL DEFAULT 0,
//  "purchase_date" date,
//  "sold_amount" numeric(10,2) NOT NULL DEFAULT 0,
//  "sold_date" date,
//  "budget_or_cost" numeric(10,2) NOT NULL DEFAULT 0,
//  "profit" numeric(10,2) NOT NULL DEFAULT 0,
//  "is_tenant_occupied" bool NOT NULL DEFAULT false,
//  "type" varchar(30) COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
//  CONSTRAINT "loan_experiences_pkey" PRIMARY KEY ("id")
// )

type LoanExperience struct {
	Base
	LoanId           int64                        `gorm:"column:loan_id;not null;default:0" json:"loan_id"`
	City             string                       `gorm:"column:city;size:20;not null;default:''" json:"city"`
	State            string                       `gorm:"column:state;size:20;not null;default:''" json:"state"`
	OwnershipOnTitle string                       `gorm:"column:ownership_on_title;size:100;not null;default:''" json:"ownership_on_title"`
	PropertyType     constants.PropertyType       `gorm:"column:property_type;size:50;not null;default:''" json:"property_type"`
	PurchaseAmount   float64                      `gorm:"column:purchase_amount;type:numeric(10,2);not null;default:0" json:"purchase_amount"`
	PurchaseDate     *sql.NullTime                `gorm:"column:purchase_date" json:"purchase_date"`
	SoldAmount       float64                      `gorm:"column:sold_amount;type:numeric(10,2);not null;default:0" json:"sold_amount"`
	SoldDate         *sql.NullTime                `gorm:"column:sold_date" json:"sold_date"`
	BudgetOrCost     float64                      `gorm:"column:budget_or_cost;type:numeric(10,2);not null;default:0" json:"budget_or_cost"`
	Profit           float64                      `gorm:"column:profit;type:numeric(10,2);not null;default:0" json:"profit"`
	IsTenantOccupied bool                         `gorm:"column:is_tenant_occupied;not null;default:false" json:"is_tenant_occupied"`
	Type             constants.LoanExperienceType `gorm:"column:type;size:30;not null;default:''" json:"type"`
	// relationships
	Loan *Loan `json:"loan" gorm:"foreignKey:LoanId;references:Id"`
}

func (m LoanExperience) TableName() string {
	return "loan_experiences"
}
