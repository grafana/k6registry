// Package k6registry contains the data model of the k6 extensions registry.
package k6registry

import _ "embed"

// Schema contains JSON schema for registry JSON.
//
//go:embed docs/registry.schema.json
var Schema []byte
