package bo

type PageResultBo[T any] struct {
	Page      int64 `json:"page"`
	Limit     int64 `json:"limit"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_page"`
	List      []T   `json:"list"`
}

type AddressBo struct {
	State    string `json:"state" copier:"State"`
	City     string `json:"city"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	ZipCode  string `json:"zip_code"`
}

type RangeInt struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type RangeFloat64 struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}
