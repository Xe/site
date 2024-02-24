package generator

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/emicklei/proto"
	"github.com/getkin/kin-openapi/openapi3"
)

const (
	googleAnyType       = "google.protobuf.Any"
	googleListValueType = "google.protobuf.ListValue"
	googleStructType    = "google.protobuf.Struct"
	googleValueType     = "google.protobuf.Value"

	googleMoneyType = "google.type.Money"
)

var (
	successDescription = "Success"
)

func (gen *generator) Handlers() []proto.Handler {
	return []proto.Handler{
		proto.WithPackage(gen.Package),
		proto.WithImport(gen.Import),
		proto.WithRPC(gen.RPC),
		proto.WithEnum(gen.Enum),
		proto.WithMessage(gen.Message),
	}
}

func (gen *generator) Package(pkg *proto.Package) {
	slog.Debug("Package handler", "package", pkg.Name)
	gen.packageName = pkg.Name
}

func (gen *generator) Import(i *proto.Import) {
	slog.Debug("Import handler", "package", gen.packageName, "filename", i.Filename)

	if _, ok := gen.importedFiles[i.Filename]; ok {
		return
	}
	gen.importedFiles[i.Filename] = struct{}{}

	// Instead of loading and generating the OpenAPI docs for the google proto definitions,
	// its known types are mapped to OpenAPI types; see aliases.go.
	if strings.Contains(i.Filename, "google/") {
		return
	}

	protoFile, err := readProtoFile(i.Filename, gen.conf.protoPaths)
	if err != nil {
		slog.Error("could not import file", "filename", i.Filename, "error", err)
		return
	}

	oldPackageName := gen.packageName

	// Override the package name for the next round of Walk calls to preserve the types full import path
	withPackage := func(pkg *proto.Package) {
		gen.packageName = pkg.Name
	}

	// additional files walked for messages and imports only
	proto.Walk(protoFile,
		proto.WithPackage(withPackage),
		proto.WithImport(gen.Import),
		proto.WithRPC(gen.RPC),
		proto.WithEnum(gen.Enum),
		proto.WithMessage(gen.Message),
	)

	gen.packageName = oldPackageName
}

func (gen *generator) RPC(rpc *proto.RPC) {
	slog.Debug("RPC handler", "package", gen.packageName, "rpc", rpc.Name, "requestType", rpc.RequestType, "returnsType", rpc.ReturnsType)

	parent, ok := rpc.Parent.(*proto.Service)
	if !ok {
		log.Panicf("parent is not proto.service")
	}
	pathName := filepath.Join("/"+gen.conf.pathPrefix+"/", gen.packageName+"."+parent.Name, rpc.Name)

	var reqMediaType *openapi3.MediaType
	switch rpc.RequestType {
	case "google.protobuf.Empty":
		reqMediaType = openapi3.NewMediaType()
	default:
		if strings.Contains(rpc.RequestType, ".") {
			reqMediaType = &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Ref: fmt.Sprintf("#/components/schemas/%s", rpc.RequestType),
				},
			}
		} else {
			reqMediaType = &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Ref: fmt.Sprintf("#/components/schemas/%s.%s", gen.packageName, rpc.RequestType),
				},
			}	
		}
	}

	var resMediaType *openapi3.MediaType
	switch rpc.ReturnsType {
	case "google.protobuf.Empty":
		resMediaType = openapi3.NewMediaType()
	default:
		if strings.Contains(rpc.ReturnsType, ".") {
			resMediaType = &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Ref: fmt.Sprintf("#/components/schemas/%s", rpc.ReturnsType),
				},
			}
		} else {
			resMediaType = &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Ref: fmt.Sprintf("#/components/schemas/%s.%s", gen.packageName, rpc.ReturnsType),
				},
			}
		}
	}

	// NOTE: Redocly does not read the "examples" (plural) field, only the "example" (singular) one.
	commentMsg, reqExamples, resExamples, err := parseComment(rpc.Comment)
	if err != nil {
		// TODO(dm): how can we surface the errors from the parser instead of panicking?
		log.Panicf("failed to parse comment %s ", err)
	}

	if len(reqExamples) > 0 {
		exampleObj := make(map[string]interface{})
		for i, example := range reqExamples {
			exampleObj[fmt.Sprintf("example %d", i)] = example
		}
		reqMediaType.Example = exampleObj
	}
	if len(resExamples) > 0 {
		exampleObj := make(map[string]interface{})
		for i, example := range resExamples {
			exampleObj[fmt.Sprintf("example %d", i)] = example
		}
		resMediaType.Example = exampleObj
	}

	gen.openAPIV3.Paths[pathName] = &openapi3.PathItem{
		Post: &openapi3.Operation{
			Description: commentMsg,
			Summary:     rpc.Name,
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: openapi3.Content{"application/json": reqMediaType},
				},
			},
			Responses: map[string]*openapi3.ResponseRef{
				"200": {
					Value: &openapi3.Response{
						Description: &successDescription,
						Content:     openapi3.Content{"application/json": resMediaType},
					},
				},
			},
		},
	}
}

func (gen *generator) Enum(enum *proto.Enum) {
	slog.Debug("Enum handler", "package", gen.packageName, "enum", enum.Name)
	values := []interface{}{}
	for _, element := range enum.Elements {
		enumField := element.(*proto.EnumField)
		values = append(values, enumField.Name)
	}

	gen.openAPIV3.Components.Schemas[gen.packageName+"."+enum.Name] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Description: description(enum.Comment),
			Type:        "string",
			Enum:        values,
		},
	}
}

func (gen *generator) Message(msg *proto.Message) {
	slog.Debug("Message handler", "package", gen.packageName, "message", msg.Name)

	schemaProps := openapi3.Schemas{}

	for _, element := range msg.Elements {
		switch val := element.(type) {
		case *proto.Message:
			//logger.logd("proto.Message")
			gen.Message(val)
		case *proto.Comment:
			//logger.logd("proto.Comment")
		case *proto.Oneof:
			//logger.logd("proto.Oneof")
		case *proto.OneOfField:
			//logger.logd("proto.OneOfField")
			gen.addField(schemaProps, val.Field, false)
		case *proto.MapField:
			//logger.logd("proto.MapField")
			gen.addField(schemaProps, val.Field, false)
		case *proto.NormalField:
			//logger.logd("proto.NormalField %q %q", val.Field.Type, val.Field.Name)
			gen.addField(schemaProps, val.Field, val.Repeated)
		default:
			slog.Error("unknown field type", "type", fmt.Sprintf("%T", element))
		}
	}

	gen.openAPIV3.Components.Schemas[gen.packageName+"."+msg.Name] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Description: description(msg.Comment),
			Type:        "object",
			Properties:  schemaProps,
		},
	}
}

func (gen *generator) addField(schemaPropsV3 openapi3.Schemas, field *proto.Field, repeated bool) {
	fieldDescription := description(field.Comment)
	fieldName := field.Name
	fieldType := field.Type
	fieldFormat := field.Type
	// map proto types to openapi
	if p, ok := typeAliases[fieldType]; ok {
		fieldType = p.Type
		fieldFormat = p.Format
	}

	if fieldType == fieldFormat {
		fieldFormat = ""
	}

	switch fieldType {
	// Build the schema for native types that don't need to reference other schemas
	// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#data-types
	case "boolean", "integer", "number", "string", "object":
		fieldSchemaV3 := openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Description: fieldDescription,
				Type:        fieldType,
				Format:      fieldFormat,
			},
		}
		if !repeated {
			schemaPropsV3[fieldName] = &fieldSchemaV3
			return
		}
		schemaPropsV3[fieldName] = &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Description: fieldDescription,
				Type:        "array",
				Format:      fieldFormat,
				Items:       &fieldSchemaV3,
			},
		}
		return

	// generate the schema for google well known complex types: https://protobuf.dev/reference/protobuf/google.protobuf/#index
	case googleAnyType:
		slog.Debug("any", "name", fieldName, "type", fieldType, "format", fieldFormat)
		gen.addGoogleAnySchema()
	case googleListValueType:
		slog.Debug("ListValue", "name", fieldName, "type", fieldType, "format", fieldFormat)
		gen.addGoogleListValueSchema()
	case googleStructType:
		slog.Debug("Struct", "name", fieldName, "type", fieldType, "format", fieldFormat)
		gen.addGoogleValueSchema() // struct depends on value
		gen.addGoogleStructSchema()
	case googleValueType:
		slog.Debug("Value", "name", fieldName, "type", fieldType, "format", fieldFormat)
		gen.addGoogleValueSchema()
	case googleMoneyType:
		slog.Debug("Money", "name", fieldName, "type", fieldType, "format", fieldFormat)
		gen.addGoogleMoneySchema()
	default:
		slog.Debug("Default", "name", fieldName, "type", fieldType, "format", fieldFormat)
	}

	// prefix custom types with the package name
	ref := fmt.Sprintf("#/components/schemas/%s", fieldType)
	if !strings.Contains(fieldType, ".") {
		ref = fmt.Sprintf("#/components/schemas/%s.%s", gen.packageName, fieldType)
	}

	if !repeated {
		schemaPropsV3[fieldName] = &openapi3.SchemaRef{
			Ref: ref,
			Value: &openapi3.Schema{
				Description: fieldDescription,
				Type:        "object",
			},
		}
		return
	}

	schemaPropsV3[fieldName] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Description: fieldDescription,
			Type:        "array",
			Items: &openapi3.SchemaRef{
				Ref: ref,
				Value: &openapi3.Schema{
					Type: "object",
				},
			},
		},
	}
}

// addGoogleAnySchema adds a schema item for the google.protobuf.Any type.
func (gen *generator) addGoogleAnySchema() {
	if _, ok := gen.openAPIV3.Components.Schemas[googleAnyType]; ok {
		return
	}
	gen.openAPIV3.Components.Schemas[googleAnyType] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Description: `
The JSON representation of an Any value uses the regular
representation of the deserialized, embedded message, with an
additional field @type which contains the type URL. Example:

	package google.profile;
	message Person {
	  string first_name = 1;
	  string last_name = 2;
	}

	{
	  "@type": "type.googleapis.com/google.profile.Person",
	  "firstName": <string>,
	  "lastName": <string>
	}

If the embedded message type is well-known and has a custom JSON
representation, that representation will be embedded adding a field
value which holds the custom JSON in addition to the @type
field. Example (for message [google.protobuf.Duration][]):

	{
	  "@type": "type.googleapis.com/google.protobuf.Duration",
	  "value": "1.212s"
	}
`,
			Type: "object",
			Properties: openapi3.Schemas{
				"@type": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "",
						Type:        "string",
						Format:      "",
					},
				},
			},
		},
	}
}

// addGoogleAnySchema adds a schema item for the google.protobuf.ListValue type.
func (gen *generator) addGoogleListValueSchema() {
	if _, ok := gen.openAPIV3.Components.Schemas[googleListValueType]; ok {
		return
	}
	gen.openAPIV3.Components.Schemas[googleListValueType] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Description: `
ListValue is a wrapper around a repeated field of values.
The JSON representation for ListValue is JSON array.
`,
			Type: "array",
			Items: &openapi3.SchemaRef{
				Value: &openapi3.Schema{
					OneOf: openapi3.SchemaRefs{
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "string",
							},
						},
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "number",
							},
						},
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "integer",
							},
						},
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "boolean",
							},
						},
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "array",
							},
						},
						&openapi3.SchemaRef{
							Value: &openapi3.Schema{
								Type: "object",
							},
						},
					},
				},
			},
		},
	}
}

func (gen *generator) addGoogleStructSchema() {
	if _, ok := gen.openAPIV3.Components.Schemas[googleStructType]; ok {
		return
	}

	gen.openAPIV3.Components.Schemas[googleStructType] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Description: `
Struct represents a structured data value, consisting of fields
which map to dynamically typed values. In some languages, 
Struct might be supported by a native representation. For example,
in scripting languages like JS a struct is represented as
an object. The details of that representation are described
together with the proto support for the language.

The JSON representation for Struct is JSON object.
`,
			Type: "object",
			Properties: openapi3.Schemas{
				"fields": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "Unordered map of dynamically typed values.",
						Type:        "object",
						AdditionalProperties: openapi3.AdditionalProperties{
							Schema: &openapi3.SchemaRef{
								Ref: "#/components/schemas/google.protobuf.Value",
							},
						},
					},
				},
			},
		},
	}
}

func (gen *generator) addGoogleValueSchema() {
	if _, ok := gen.openAPIV3.Components.Schemas[googleValueType]; ok {
		return
	}

	gen.openAPIV3.Components.Schemas[googleValueType] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Description: `
Value represents a dynamically typed value which can be either
null, a number, a string, a boolean, a recursive struct value, or a
list of values. A producer of value is expected to set one of that
variants, absence of any variant indicates an error.
				
The JSON representation for Value is JSON value.
`,
			OneOf: openapi3.SchemaRefs{
				&openapi3.SchemaRef{Value: &openapi3.Schema{Type: "string"}},
				&openapi3.SchemaRef{Value: &openapi3.Schema{Type: "number"}},
				&openapi3.SchemaRef{Value: &openapi3.Schema{Type: "integer"}},
				&openapi3.SchemaRef{Value: &openapi3.Schema{Type: "boolean"}},
				&openapi3.SchemaRef{Ref: "#/components/schemas/google.protobuf.Struct"},
				&openapi3.SchemaRef{Ref: "#/components/schemas/google.protobuf.ListValue"},
			},
		},
	}
}

func (gen *generator) addGoogleMoneySchema() {
	if _, ok := gen.openAPIV3.Components.Schemas[googleMoneyType]; ok {
		return
	}

	gen.openAPIV3.Components.Schemas[googleMoneyType] = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			Description: `Represents an amount of money with its currency type`,
			Type:        "object",
			Properties: openapi3.Schemas{
				"currency_code": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "The 3-letter currency code defined in ISO 4217.",
						Type:        "string",
					},
				},
				"units": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "The whole units of the amount.\nFor example if `currencyCode` is `\"USD\"`, then 1 unit is one US dollar.",
						Type:        "integer",
						Format:      "int64",
					},
				},
				"nanos": &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Description: "Number of nano (10^-9) units of the amount.\nThe value must be between -999,999,999 and +999,999,999 inclusive.\nIf `units` is positive, `nanos` must be positive or zero.\nIf `units` is zero, `nanos` can be positive, zero, or negative.\nIf `units` is negative, `nanos` must be negative or zero.\nFor example $-1.75 is represented as `units`=-1 and `nanos`=-750,000,000.",
						Type:        "integer",
						Format:      "int32",
					},
				},
			},
		},
	}
}

func description(comment *proto.Comment) string {
	if comment == nil {
		return ""
	}
	result := []string{}
	for _, line := range comment.Lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// parseComment parses the comment for an RPC method and returns the description, request examples, and response examples.
// it looks for the labels req-example: and res-example: to extract the JSON payload samples.
func parseComment(comment *proto.Comment) (string, []map[string]interface{}, []map[string]interface{}, error) {
	if comment == nil {
		return "", nil, nil, nil
	}
	reqExamples := []map[string]interface{}{}
	respExamples := []map[string]interface{}{}
	message := ""
	for _, line := range comment.Lines {
		line = strings.TrimLeftFunc(line, unicode.IsSpace)
		if strings.HasPrefix(line, "req-example:") {
			parts := strings.Split(line, "req-example:")
			example := map[string]interface{}{}
			if err := json.Unmarshal([]byte(parts[1]), &example); err != nil {
				return "", nil, nil, fmt.Errorf("failed to parse req-example %q: %v", parts[1], err)
			}
			reqExamples = append(reqExamples, example)
		} else if strings.HasPrefix(line, "res-example:") {
			parts := strings.Split(line, "res-example:")
			example := map[string]interface{}{}
			if err := json.Unmarshal([]byte(parts[1]), &example); err != nil {
				return "", nil, nil, fmt.Errorf("failed to parse res-example %q: %v", parts[1], err)
			}
			respExamples = append(respExamples, example)
		} else {
			message = fmt.Sprintf("%s\n%s", message, line)
		}
	}
	return message, reqExamples, respExamples, nil
}
