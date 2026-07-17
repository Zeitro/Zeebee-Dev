package config

type LlmConf struct {
	Providers []LlmProvider
}

type LlmProvider struct {
	Name    string
	ApiKey  string
	BaseUrl string
}

type LlmBizConf struct {
	Provider string
	Model    string
	// MaxTokens overrides the provider default output token cap (claude only).
	// 0 leaves the client default (4096) in place.
	MaxTokens int64 `json:",optional"`
	// TimeoutSeconds overrides the per-request HTTP timeout. 0 leaves the
	// client default (60s) in place.
	TimeoutSeconds int `json:",optional"`
}
