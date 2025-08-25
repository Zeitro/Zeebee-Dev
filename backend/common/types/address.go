package types

type Address struct {
	State    string `json:"state" copier:"State"`
	City     string `json:"city"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	ZipCode  string `json:"zip_code"`
}
