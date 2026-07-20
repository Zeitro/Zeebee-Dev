package config

type LlmConf struct {
	Providers []LlmProvider `json:",optional"`
}

type LlmProvider struct {
	Name    string `json:",optional"`
	ApiKey  string `json:",optional"`
	BaseUrl string `json:",optional"`
}

type LlmBizConf struct {
	Provider string `json:",optional"`
	Model    string `json:",optional"`
	// MaxTokens overrides the provider default output token cap (claude only).
	// 0 leaves the client default (4096) in place.
	MaxTokens int64 `json:",optional"`
	// TimeoutSeconds overrides the per-request HTTP timeout. 0 leaves the
	// client default (60s) in place.
	TimeoutSeconds int `json:",optional"`
}
