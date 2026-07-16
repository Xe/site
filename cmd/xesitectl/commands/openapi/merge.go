// Package openapi merges the per-proto-package OpenAPI documents that
// protoc-gen-connect-openapi writes into gen/ into a single document
// describing all of xeiaso.net's API surface.
package openapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"slices"
)

// A Fragment is one generated OpenAPI document awaiting merge. Path is carried
// along so conflicts can name the file that introduced them.
type Fragment struct {
	Path string
	Data []byte
}

// Merge folds each fragment's paths, schemas and tags into base and returns the
// combined document.
//
// Fragments routinely repeat definitions: every Twirp service embeds
// TwirpError, and shared protos like google.protobuf.Timestamp appear in most
// of them. An identical redefinition is therefore expected and ignored. Only a
// name bound to two different definitions is an error, since merge order would
// otherwise silently decide which one wins.
func Merge(base []byte, fragments []Fragment) ([]byte, error) {
	doc := map[string]any{}
	if err := decodeJSON(base, &doc); err != nil {
		return nil, fmt.Errorf("openapi: parse base document: %w", err)
	}

	paths, err := objectAt(doc, "paths")
	if err != nil {
		return nil, fmt.Errorf("openapi: base document: %w", err)
	}

	components, err := objectAt(doc, "components")
	if err != nil {
		return nil, fmt.Errorf("openapi: base document: %w", err)
	}

	schemas, err := objectAt(components, "schemas")
	if err != nil {
		return nil, fmt.Errorf("openapi: base document: components: %w", err)
	}

	// Tags are a list in the document but dedupe by name, so they get indexed
	// into a map for the duration of the merge and re-listed at the end.
	baseTags, err := fragList(doc, "tags")
	if err != nil {
		return nil, fmt.Errorf("openapi: base document: %w", err)
	}

	tags := map[string]any{}
	if err := indexTags(tags, baseTags, "base document"); err != nil {
		return nil, err
	}

	for _, f := range fragments {
		fragDoc := map[string]any{}
		if err := decodeJSON(f.Data, &fragDoc); err != nil {
			return nil, fmt.Errorf("openapi: parse %s: %w", f.Path, err)
		}

		fp, ok, err := fragObject(fragDoc, "paths")
		if err != nil {
			return nil, fmt.Errorf("openapi: %s: %w", f.Path, err)
		}
		if ok {
			if err := mergeInto(paths, fp, "path", f.Path); err != nil {
				return nil, err
			}
		}

		fc, ok, err := fragObject(fragDoc, "components")
		if err != nil {
			return nil, fmt.Errorf("openapi: %s: %w", f.Path, err)
		}
		if ok {
			fs, ok, err := fragObject(fc, "schemas")
			if err != nil {
				return nil, fmt.Errorf("openapi: %s: components: %w", f.Path, err)
			}
			if ok {
				if err := mergeInto(schemas, fs, "schema", f.Path); err != nil {
					return nil, err
				}
			}
		}

		// Message-only protos have no services, so connect-openapi omits their
		// tags entirely. An absent or null tags key is normal rather than a
		// defect; only a tags value of the wrong shape is an error.
		fragTags, err := fragList(fragDoc, "tags")
		if err != nil {
			return nil, fmt.Errorf("openapi: %s: %w", f.Path, err)
		}
		if err := indexTags(tags, fragTags, f.Path); err != nil {
			return nil, err
		}
	}

	doc["tags"] = valuesByKey(tags)

	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("openapi: encode merged document: %w", err)
	}

	return append(out, '\n'), nil
}

// mergeInto copies src's entries into dst, rejecting any key that src and dst
// define differently. kind and srcPath only shape the error message.
func mergeInto(dst, src map[string]any, kind, srcPath string) error {
	for k, v := range src {
		if old, ok := dst[k]; ok {
			if !reflect.DeepEqual(old, v) {
				return fmt.Errorf("openapi: %s: conflicting %s %q: already defined with a different shape", srcPath, kind, k)
			}

			continue
		}

		dst[k] = v
	}

	return nil
}

// indexTags folds an OpenAPI tag list into dst, keyed by tag name.
func indexTags(dst map[string]any, tags []any, srcPath string) error {
	for _, t := range tags {
		tag, ok := t.(map[string]any)
		if !ok {
			return fmt.Errorf("openapi: %s: tag is not an object", srcPath)
		}

		name, ok := tag["name"].(string)
		if !ok {
			return fmt.Errorf("openapi: %s: tag name is missing or not a string", srcPath)
		}

		if err := mergeInto(dst, map[string]any{name: tag}, "tag", srcPath); err != nil {
			return err
		}
	}

	return nil
}

// objectAt returns doc[key] as an object, creating it when absent so the base
// document doesn't have to spell out every empty container.
func objectAt(doc map[string]any, key string) (map[string]any, error) {
	v, ok := doc[key]
	if !ok || v == nil {
		m := map[string]any{}
		doc[key] = m

		return m, nil
	}

	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%q is not an object", key)
	}

	return m, nil
}

// decodeJSON parses data into v using json.Number for numeric literals
// instead of float64, so integers wider than 2^53 (an int64 maximum or
// default value, for instance) survive a decode/re-encode round trip
// byte-exact instead of being silently rounded. It otherwise matches
// json.Unmarshal, including rejecting trailing non-whitespace content.
func decodeJSON(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(v); err != nil {
		return err
	}

	if _, err := dec.Token(); err != io.EOF {
		return fmt.Errorf("unexpected trailing data after JSON value")
	}

	return nil
}

// fragObject returns m[key] as an object. A key that is absent or explicitly
// null is normal for a fragment (e.g. message-only protos omit "tags"
// entirely) and reports as not-ok without an error. A key that is present
// with a different, non-object shape is a defect in the generator output and
// must not be silently dropped, so it is reported as an error.
func fragObject(m map[string]any, key string) (obj map[string]any, ok bool, err error) {
	v, present := m[key]
	if !present || v == nil {
		return nil, false, nil
	}

	obj, ok = v.(map[string]any)
	if !ok {
		return nil, false, fmt.Errorf("%q is not an object", key)
	}

	return obj, true, nil
}

// fragList returns m[key] as a JSON list, with the same absent/null-is-fine,
// wrong-type-is-an-error rules as fragObject.
func fragList(m map[string]any, key string) ([]any, error) {
	v, present := m[key]
	if !present || v == nil {
		return nil, nil
	}

	l, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("%q is not a list", key)
	}

	return l, nil
}

// valuesByKey returns m's values ordered by key. Map iteration order is random
// and gen/ is committed, so an unsorted list would produce a spurious diff on
// every run.
func valuesByKey(m map[string]any) []any {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	out := make([]any, 0, len(keys))
	for _, k := range keys {
		out = append(out, m[k])
	}

	return out
}
