package static

import _ "embed"

//go:embed api.json
var ApiStatic []byte

//go:embed resource.json
var ResourceStatic []byte
