// Package swagger provides embedded swagger documentation for the borrower API.
package swagger

import "embed"

// SwaggerJSON contains the embedded swagger documentation file.
//
//go:embed borrower.json
var SwaggerJSON embed.FS
