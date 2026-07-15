// Package openapi merges the per-proto-package OpenAPI documents that
// protoc-gen-connect-openapi writes into gen/ into a single document
// describing all of xeiaso.net's API surface.
package openapi

import (
	"encoding/json"
	"fmt"
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
	if err := json.Unmarshal(base, &doc); err != nil {
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
	tags := map[string]any{}
	if err := indexTags(tags, list(doc["tags"]), "base document"); err != nil {
		return nil, err
	}

	for _, f := range fragments {
		frag := map[string]any{}
		if err := json.Unmarshal(f.Data, &frag); err != nil {
			return nil, fmt.Errorf("openapi: parse %s: %w", f.Path, err)
		}

		if fp, ok := frag["paths"].(map[string]any); ok {
			if err := mergeInto(paths, fp, "path", f.Path); err != nil {
				return nil, err
			}
		}

		if fc, ok := frag["components"].(map[string]any); ok {
			if fs, ok := fc["schemas"].(map[string]any); ok {
				if err := mergeInto(schemas, fs, "schema", f.Path); err != nil {
					return nil, err
				}
			}
		}

		// Message-only protos have no services, so connect-openapi omits their
		// tags entirely. A nil tags key is normal rather than a defect.
		if err := indexTags(tags, list(frag["tags"]), f.Path); err != nil {
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
			return fmt.Errorf("openapi: %s: tag has no name", srcPath)
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

// list returns v as a JSON list, or nil if it isn't one.
func list(v any) []any {
	l, _ := v.([]any)

	return l
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
