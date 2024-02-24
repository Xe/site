package generator

var typeAliases = map[string]struct {
	Type, Format string
}{
	// proto numeric types
	"int32":    {Type: "integer", Format: "int32"},
	"uint32":   {Type: "integer", Format: "uint32"},
	"sint32":   {Type: "integer", Format: "int32"},
	"fixed32":  {Type: "integer", Format: "int32"},
	"sfixed32": {Type: "integer", Format: "int32"},

	// proto numeric types, 64bit
	"int64":    {Type: "string", Format: "int64"},
	"uint64":   {Type: "string", Format: "uint64"},
	"sint64":   {Type: "string", Format: "int64"},
	"fixed64":  {Type: "string", Format: "int64"},
	"sfixed64": {Type: "string", Format: "int64"},

	"double": {Type: "number", Format: "double"},
	"float":  {Type: "number", Format: "float"},

	// effectively copies google.protobuf.BytesValue
	"bytes": {
		Type:   "string",
		Format: "byte",
	},

	// It is what it is
	"bool": {
		Type:   "boolean",
		Format: "boolean",
	},

	"google.protobuf.Timestamp": {
		Type:   "string",
		Format: "date-time",
	},
	"google.protobuf.Duration": {
		Type: "string",
	},
	"google.protobuf.StringValue": {
		Type: "string",
	},
	"google.protobuf.BytesValue": {
		Type:   "string",
		Format: "byte",
	},
	"google.protobuf.Int32Value": {
		Type:   "integer",
		Format: "int32",
	},
	"google.protobuf.UInt32Value": {
		Type:   "integer",
		Format: "uint32",
	},
	"google.protobuf.Int64Value": {
		Type:   "string",
		Format: "int64",
	},
	"google.protobuf.UInt64Value": {
		Type:   "string",
		Format: "uint64",
	},
	"google.protobuf.FloatValue": {
		Type:   "number",
		Format: "float",
	},
	"google.protobuf.DoubleValue": {
		Type:   "number",
		Format: "double",
	},
	"google.protobuf.BoolValue": {
		Type:   "boolean",
		Format: "boolean",
	},
	"google.protobuf.Empty": {
		Type: "object",
	},

	"google.type.DateTime": {
		Type:   "string",
		Format: "date-time",
	},
}
