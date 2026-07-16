package config

type LlmConf struct {
	Providers []LlmProvider
}

type LlmProvider struct {
	Name string
	ApiKey string
	BaseUrl string
}


type LlmBizConf struct {
	Provider string
	Model string
}