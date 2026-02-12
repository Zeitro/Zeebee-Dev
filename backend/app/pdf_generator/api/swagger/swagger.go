// Package swagger provides embedded swagger documentation for the lender API.
package swagger

import "embed"

// SwaggerJSON contains the embedded swagger documentation file.
//
//go:embed pdf_generator.json
var SwaggerJSON embed.FS
