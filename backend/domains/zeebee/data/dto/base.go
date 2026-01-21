package dto

type AddressDto struct {
	State    string `json:"state" copier:"State"`
	City     string `json:"city"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	ZipCode  string `json:"zip_code"`
}

type BasePageRequestDto struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type BasePageFilter struct {
	Keyword string `json:"keyword"`
}

type BasePageSorter struct {
	Column string `json:"column"`
	Sort   string `json:"sort"`
}

type RangeInt struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type RangeFloat64 struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}
