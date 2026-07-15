package openapi

import (
	"encoding/json"
	"strings"
	"testing"
)

// minimalBase is the smallest base document Merge accepts, kept separate from
// the real base.json so these tests don't fail when the metadata changes.
const minimalBase = `{
  "openapi": "3.1.0",
  "info": {"title": "test", "contact": {"email": "me@xeiaso.net"}},
  "security": [],
  "tags": [],
  "paths": {},
  "components": {"schemas": {}}
}`

func frag(path, body string) Fragment {
	return Fragment{Path: path, Data: []byte(body)}
}

func TestMerge(t *testing.T) {
	for _, tt := range []struct {
		name      string
		base      string
		fragments []Fragment
		wantErr   string
		check     func(t *testing.T, doc map[string]any)
	}{
		{
			name: "combines paths from every fragment",
			base: minimalBase,
			fragments: []Fragment{
				frag("a.openapi.json", `{"paths": {"/twirp/a.A/Get": {"post": {}}}}`),
				frag("b.openapi.json", `{"paths": {"/twirp/b.B/Get": {"post": {}}}}`),
			},
			check: func(t *testing.T, doc map[string]any) {
				paths := doc["paths"].(map[string]any)
				if len(paths) != 2 {
					t.Errorf("got %d paths, want 2: %v", len(paths), paths)
				}
				if _, ok := paths["/twirp/a.A/Get"]; !ok {
					t.Error("missing path from fragment a")
				}
				if _, ok := paths["/twirp/b.B/Get"]; !ok {
					t.Error("missing path from fragment b")
				}
			},
		},
		{
			name: "identical duplicate schema is not an error",
			base: minimalBase,
			fragments: []Fragment{
				frag("a.openapi.json", `{"components": {"schemas": {"TwirpError": {"type": "object"}}}}`),
				frag("b.openapi.json", `{"components": {"schemas": {"TwirpError": {"type": "object"}}}}`),
			},
			check: func(t *testing.T, doc map[string]any) {
				schemas := doc["components"].(map[string]any)["schemas"].(map[string]any)
				if len(schemas) != 1 {
					t.Errorf("got %d schemas, want 1: %v", len(schemas), schemas)
				}
			},
		},
		{
			name: "conflicting schema is an error",
			base: minimalBase,
			fragments: []Fragment{
				frag("a.openapi.json", `{"components": {"schemas": {"Thing": {"type": "object"}}}}`),
				frag("b.openapi.json", `{"components": {"schemas": {"Thing": {"type": "string"}}}}`),
			},
			wantErr: `conflicting schema "Thing"`,
		},
		{
			name: "conflicting path is an error",
			base: minimalBase,
			fragments: []Fragment{
				frag("a.openapi.json", `{"paths": {"/twirp/a.A/Get": {"post": {"operationId": "one"}}}}`),
				frag("b.openapi.json", `{"paths": {"/twirp/a.A/Get": {"post": {"operationId": "two"}}}}`),
			},
			wantErr: `conflicting path "/twirp/a.A/Get"`,
		},
		{
			name: "fragment with null tags is fine",
			base: minimalBase,
			fragments: []Fragment{
				frag("protofeed.openapi.json", `{"tags": null, "components": {"schemas": {"protofeed.Item": {"type": "object"}}}}`),
			},
			check: func(t *testing.T, doc map[string]any) {
				if got := doc["tags"].([]any); len(got) != 0 {
					t.Errorf("got %d tags, want 0", len(got))
				}
			},
		},
		{
			name: "fragment with no fields at all is fine",
			base: minimalBase,
			fragments: []Fragment{
				frag("empty.openapi.json", `{}`),
			},
			check: func(t *testing.T, doc map[string]any) {
				if got := doc["tags"].([]any); len(got) != 0 {
					t.Errorf("got %d tags, want 0", len(got))
				}
				paths := doc["paths"].(map[string]any)
				if len(paths) != 0 {
					t.Errorf("got %d paths, want 0: %v", len(paths), paths)
				}
			},
		},
		{
			name: "fragment with paths of the wrong type is an error",
			base: minimalBase,
			fragments: []Fragment{
				frag("bad.openapi.json", `{"paths": "not an object"}`),
			},
			wantErr: `bad.openapi.json: "paths" is not an object`,
		},
		{
			name: "fragment with components of the wrong type is an error",
			base: minimalBase,
			fragments: []Fragment{
				frag("bad.openapi.json", `{"components": "not an object"}`),
			},
			wantErr: `bad.openapi.json: "components" is not an object`,
		},
		{
			name: "fragment with components.schemas of the wrong type is an error",
			base: minimalBase,
			fragments: []Fragment{
				frag("bad.openapi.json", `{"components": {"schemas": "not an object"}}`),
			},
			wantErr: `bad.openapi.json: components: "schemas" is not an object`,
		},
		{
			name: "fragment with tags of the wrong type is an error",
			base: minimalBase,
			fragments: []Fragment{
				frag("bad.openapi.json", `{"tags": "not a list"}`),
			},
			wantErr: `bad.openapi.json: "tags" is not a list`,
		},
		{
			name: "tags are deduped and sorted by name",
			base: minimalBase,
			fragments: []Fragment{
				frag("b.openapi.json", `{"tags": [{"name": "zeta"}, {"name": "gamma"}]}`),
				frag("a.openapi.json", `{"tags": [{"name": "alpha"}]}`),
			},
			check: func(t *testing.T, doc map[string]any) {
				tags := doc["tags"].([]any)
				if len(tags) != 3 {
					t.Fatalf("got %d tags, want 3: %v", len(tags), tags)
				}
				var names []string
				for _, tag := range tags {
					names = append(names, tag.(map[string]any)["name"].(string))
				}
				want := []string{"alpha", "gamma", "zeta"}
				for i, name := range names {
					if name != want[i] {
						t.Errorf("got tag order %v, want %v", names, want)
						break
					}
				}
			},
		},
		{
			name: "conflicting tag is an error",
			base: minimalBase,
			fragments: []Fragment{
				frag("a.openapi.json", `{"tags": [{"name": "svc", "description": "one"}]}`),
				frag("b.openapi.json", `{"tags": [{"name": "svc", "description": "two"}]}`),
			},
			wantErr: `conflicting tag "svc"`,
		},
		{
			name:      "base metadata survives the merge",
			base:      minimalBase,
			fragments: []Fragment{frag("a.openapi.json", `{"paths": {"/twirp/a.A/Get": {"post": {}}}}`)},
			check: func(t *testing.T, doc map[string]any) {
				if got := doc["openapi"]; got != "3.1.0" {
					t.Errorf("got openapi %v, want 3.1.0", got)
				}
				info := doc["info"].(map[string]any)
				if got := info["title"]; got != "test" {
					t.Errorf("got title %v, want test", got)
				}
				if got := info["contact"].(map[string]any)["email"]; got != "me@xeiaso.net" {
					t.Errorf("got contact email %v, want me@xeiaso.net", got)
				}
				sec, ok := doc["security"].([]any)
				if !ok || len(sec) != 0 {
					t.Errorf("got security %v, want an empty list", doc["security"])
				}
			},
		},
		{
			name:      "invalid fragment json names the file",
			base:      minimalBase,
			fragments: []Fragment{frag("broken.openapi.json", `{not json`)},
			wantErr:   "broken.openapi.json",
		},
		{
			name:      "invalid base json is an error",
			base:      `{not json`,
			fragments: nil,
			wantErr:   "parse base document",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Merge([]byte(tt.base), tt.fragments)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("got nil error, want one containing %q", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("got error %q, want it to contain %q", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Merge: %v", err)
			}

			doc := map[string]any{}
			if err := json.Unmarshal(got, &doc); err != nil {
				t.Fatalf("merged output is not valid JSON: %v", err)
			}

			if tt.check != nil {
				tt.check(t, doc)
			}
		})
	}
}

func TestMergeOutputFormat(t *testing.T) {
	got, err := Merge([]byte(minimalBase), []Fragment{
		frag("a.openapi.json", `{"paths": {"/twirp/a.A/Get": {"post": {}}}}`),
	})
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}

	if len(got) == 0 || got[len(got)-1] != '\n' {
		t.Fatalf("output does not end with a newline: %q", got)
	}
	if strings.HasSuffix(string(got), "\n\n") {
		t.Fatalf("output has more than one trailing newline: %q", got)
	}

	// json.MarshalIndent with a two-space prefix nests the top-level
	// "paths" key's first entry two spaces deep and its own children four.
	if !strings.Contains(string(got), "\n  \"paths\": {\n    \"/twirp/a.A/Get\": {\n") {
		t.Fatalf("output is not indented with two spaces per level:\n%s", got)
	}
}

func TestMergeIsDeterministic(t *testing.T) {
	fragments := []Fragment{
		frag("a.openapi.json", `{"tags": [{"name": "zeta"}], "paths": {"/z": {}}, "components": {"schemas": {"Z": {"type": "object"}}}}`),
		frag("b.openapi.json", `{"tags": [{"name": "alpha"}], "paths": {"/a": {}}, "components": {"schemas": {"A": {"type": "object"}}}}`),
	}

	first, err := Merge([]byte(minimalBase), fragments)
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}

	// gen/ is committed, so unstable output would dirty the tree on every run.
	for i := range 10 {
		next, err := Merge([]byte(minimalBase), fragments)
		if err != nil {
			t.Fatalf("Merge run %d: %v", i, err)
		}
		if string(next) != string(first) {
			t.Fatalf("run %d differs from run 0", i)
		}
	}
}
