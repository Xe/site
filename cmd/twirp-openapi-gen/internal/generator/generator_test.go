package generator

import (
	"flag"
	"log/slog"
	"os"
	"strings"
	"testing"
)

type ProtoRPC struct {
	name   string
	input  string
	output string
	desc   string
}

type ProtoMessage struct {
	name   string
	fields []ProtoField
}

type ProtoField struct {
	name      string
	fieldType string
	format    string
	desc      string
	enums     []string
	ref       string
	itemsRef  string
	itemsType string
}

var (
	verbose = flag.Bool("slog.verbose", false, "print debug logs to the console")
)

func init() {
	if *verbose {
		h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
		slog.SetDefault(slog.New(h))
	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
		})))
	}
}

func TestGenerator(t *testing.T) {
	flag.Parse()

	opts := []Option{
		ProtoPaths([]string{"./testdata/paymentapis", "./testdata/petapis"}),
		Servers([]string{"https://example.com"}),
		Title("Test"),
		DocVersion("0.1"),
		Format("json"),
	}
	gen, err := NewGenerator([]string{"./testdata/petapis/pet/v1/pet.proto"}, opts...)
	if err != nil {
		t.Fatal(err)
	}
	openAPI, err := gen.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if err := gen.Save("./testdata/doc.json"); err != nil {
		t.Fatal(err)
	}

	pkgName := "pet.v1"
	serviceName := "PetStoreService"
	rpcs := []ProtoRPC{
		{
			name:   "GetPet",
			input:  "GetPetRequest",
			output: "GetPetResponse",
		},
	}
	messages := []ProtoMessage{
		{
			name: "GetPetRequest",
			fields: []ProtoField{
				{
					name:      "pet_id",
					fieldType: "string",
				},
			},
		},
		{
			name: "Pet",
			fields: []ProtoField{
				{
					name:      "pet_type",
					fieldType: "object",
					ref:       "#/components/schemas/pet.v1.PetType",
					enums: []string{
						"PET_TYPE_UNSPECIFIED",
						"PET_TYPE_CAT",
						"PET_TYPE_DOG",
						"PET_TYPE_SNAKE",
						"PET_TYPE_HAMSTER",
					},
				},
				{
					name:      "pet_types",
					fieldType: "array",
					itemsRef:  "#/components/schemas/pet.v1.PetType",
				},
				{
					name:      "tags",
					fieldType: "array",
					itemsType: "string",
				},
				{
					name:      "pet_id",
					fieldType: "string",
					desc:      "pet_id is an auto-generated id for the pet\\nthe id uniquely identifies a pet in the system",
				},
				{
					name:      "name",
					fieldType: "string",
				},
				{
					name:      "created_at",
					fieldType: "string",
					format:    "date-time",
				},
				{
					name:      "vet",
					fieldType: "object",
					ref:       "#/components/schemas/pet.v1.Vet",
				},
				{
					name:      "vets",
					fieldType: "array",
					itemsRef:  "#/components/schemas/pet.v1.Vet",
					itemsType: "object",
				},
			},
		},
	}

	t.Run("RPC", func(t *testing.T) {
		for _, rpc := range rpcs {
			pathName := "/" + pkgName + "." + serviceName + "/" + rpc.name
			path, ok := openAPI.Paths[pathName]
			if !ok {
				t.Errorf("%s: missing rpc %q", pathName, rpc.name)
			}

			if path.Description != rpc.desc {
				t.Errorf("%s: expected desc %q but got %q", pathName, rpc.desc, path.Description)
			}

			post := path.Post
			if post == nil {
				t.Errorf("%s: missing post", pathName)
				continue
			}

			if post.Summary != rpc.name {
				t.Errorf("%s: expected summary %q but got %q", pathName, rpc.name, post.Summary)
			}

			requestBodyRef := post.RequestBody
			if requestBodyRef == nil {
				t.Errorf("%s: missing request body", pathName)
				continue
			}

			// request
			{
				requestBody := requestBodyRef.Value
				if requestBody == nil {
					t.Errorf("%s: missing request body", pathName)
					continue
				}

				mediaType, ok := requestBody.Content["application/json"]
				if !ok {
					t.Errorf("%s: missing content type", pathName)
					continue
				}

				if mediaType.Schema == nil {
					t.Errorf("%s: missing media type schema", pathName)
					continue
				}

				expectedRef := "#/components/schemas/" + pkgName + "." + rpc.input
				if mediaType.Schema.Ref != expectedRef {
					t.Errorf("%s: expected ref %q but got %q", pathName, expectedRef, mediaType.Schema.Ref)
				}
			}

			// response
			{
				respRef := post.Responses["200"]
				if respRef == nil {
					t.Errorf("%s: missing resp", pathName)
					continue
				}

				resp := respRef.Value
				if resp == nil {
					t.Errorf("%s: missing resp", pathName)
					continue
				}

				mediaType, ok := resp.Content["application/json"]
				if !ok {
					t.Errorf("%s: missing content type", pathName)
					continue
				}

				if mediaType.Schema == nil {
					t.Errorf("%s: missing media type schema", pathName)
					continue
				}

				expectedRef := "#/components/schemas/" + pkgName + "." + rpc.output
				if mediaType.Schema.Ref != expectedRef {
					t.Errorf("%s: expected ref %q but got %q", pathName, expectedRef, mediaType.Schema.Ref)
				}
			}
		}
	})

	t.Run("Messages", func(*testing.T) {
		for _, message := range messages {
			schemaName := "" + pkgName + "." + message.name
			schema, ok := openAPI.Components.Schemas[schemaName]
			if !ok {
				t.Errorf("%s: missing message %q", schemaName, message.name)
			}
			if schema.Value == nil {
				t.Errorf("%s: missing component", schemaName)
				continue
			}
			properties := schema.Value.Properties
			for _, messageField := range message.fields {
				propertyRef, ok := properties[messageField.name]
				if !ok {
					t.Errorf("%s: missing property %q", schemaName, messageField.name)
				}

				if propertyRef == nil || propertyRef.Value == nil {
					t.Errorf("%s: missing property ref", schemaName)
					continue
				}

				property := propertyRef.Value
				if property.Type != messageField.fieldType {
					t.Errorf("%s: %q expected property type %q but got %q", schemaName, message.name, messageField.fieldType, property.Type)
					continue
				}

				if messageField.format != "" {
					if messageField.format != "" && property.Format != messageField.format {
						t.Errorf("%s: expected property format %q but got %q", schemaName, messageField.format, property.Format)
						continue
					}
				}

				if propertyRef.Ref != messageField.ref {
					t.Errorf("%s: %q expected reference %q but got %q", schemaName, messageField.name, messageField.ref, propertyRef.Ref)
				}

				// check the reference schema
				if messageField.ref != "" {
					refParts := strings.Split(messageField.ref, "/")
					// the reference schema has the format of #/components/schemas/<type> so we need to get the last part
					schemaRef, ok := openAPI.Components.Schemas[refParts[len(refParts)-1]]
					if !ok {
						t.Errorf("%s: %q expected reference schema %q but got nil", schemaName, messageField.name, messageField.ref)
					} else {
						// check if the schema reference has the expected enum values
						if len(messageField.enums) > 0 {
							if schemaRef.Value.Enum == nil {
								t.Errorf("%s: %q expected reference schema enums %q but got nil", schemaName, messageField.name, messageField.ref)
							} else {
								enums := map[string]struct{}{}
								for _, e := range schemaRef.Value.Enum {
									enums[e.(string)] = struct{}{}
								}
								for _, e := range messageField.enums {
									if _, ok := enums[e]; !ok {
										t.Errorf("%s: %q expected reference schema enum %q to have %q but got nil", schemaName, messageField.name, messageField.ref, e)
									}
								}
							}
						}
					}
				}

				if property.Type == "array" {
					if property.Items == nil || property.Items.Value == nil {
						t.Errorf("%s: missing property enum array items", schemaName)
					}
					// only check the array items type if it's not a reference
					if messageField.itemsRef == "" && (property.Items.Value.Type != messageField.itemsType) {
						t.Errorf("%s: expected %s items type %q but got %q", schemaName, messageField.name, messageField.itemsType, property.Items.Value.Type)
					}
					// check the array items reference schema
					if property.Items.Ref != messageField.itemsRef {
						t.Errorf("%s: expected %s items ref %q but got %q", schemaName, messageField.name, messageField.itemsRef, property.Items.Ref)
					}
				}
			}
		}
	})
}
