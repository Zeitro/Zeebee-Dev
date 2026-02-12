package enums

type EnumOpt struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

var enumMap map[string][]EnumOpt = map[string][]EnumOpt{}

func Register(name string, opts []EnumOpt) {
	enumMap[name] = opts
}

func Get(name string) []EnumOpt {
	return enumMap[name]
}
