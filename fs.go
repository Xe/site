package xeiaso

import "embed"

//go:embed bin/*.wasm
var Bin embed.FS

//go:embed data/toots/*.json data/users/*.json
var Data embed.FS

//go:embed blog/*.markdown gallery/*.markdown talks/*.markdown
var Markdown embed.FS

//go:embed tmpl/*.html
var Templates embed.FS

//go:embed static/*
var Static embed.FS
