package openapi

import _ "embed"

// baseDoc is the xesite-specific metadata that generated fragments merge into:
// the title, contact details, server list, and the empty security list that
// declares these endpoints take no authentication. The generator emits none of
// this, so it lives here and is layered on after the fact.
//
//go:embed base.json
var baseDoc []byte
