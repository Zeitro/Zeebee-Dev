package gormx

type WhereRaw struct {
	Raw  string
	Args []interface{}
}

// QueryOption query option
type QueryOption struct {
	Filter   []WhereRaw
	Order    []string
	Fields   []string
	Preloads []string
}

type QueryOptionFunc func(opt *QueryOption)

// NewQueryOption new query option
func NewQueryOption(opts ...QueryOptionFunc) QueryOption {
	opt := &QueryOption{}
	for _, o := range opts {
		o(opt)
	}
	return *opt
}

// WithFilterOpt with filter option
func WithFilterOpt(filter []WhereRaw) QueryOptionFunc {
	return func(opt *QueryOption) {
		opt.Filter = filter
	}
}

// WithOrderOpt with order option
func WithOrderOpt(order []string) QueryOptionFunc {
	return func(opt *QueryOption) {
		opt.Order = order
	}
}

// WithFieldsOpt with fields option
func WithFieldsOpt(fields []string) QueryOptionFunc {
	return func(opt *QueryOption) {
		opt.Fields = fields
	}
}

func WithPreloadsOpt(preloads []string) QueryOptionFunc {
	return func(opt *QueryOption) {
		opt.Preloads = preloads
	}
}
