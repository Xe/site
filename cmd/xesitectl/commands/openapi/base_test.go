package openapi

import (
	"encoding/json"
	"testing"
)

// TestBaseDocument guards the embedded metadata against typos: it is hand-
// written JSON that nothing else parses until a release runs.
func TestBaseDocument(t *testing.T) {
	doc := map[string]any{}
	if err := json.Unmarshal(baseDoc, &doc); err != nil {
		t.Fatalf("base.json is not valid JSON: %v", err)
	}

	if got := doc["openapi"]; got != "3.1.0" {
		t.Errorf("got openapi %v, want 3.1.0", got)
	}

	info, ok := doc["info"].(map[string]any)
	if !ok {
		t.Fatal("base.json has no info object")
	}

	if got := info["title"]; got != "xeiaso.net" {
		t.Errorf("got info.title %v, want xeiaso.net", got)
	}

	contact, ok := info["contact"].(map[string]any)
	if !ok {
		t.Fatal("base.json has no info.contact object")
	}

	if got := contact["email"]; got != "me@xeiaso.net" {
		t.Errorf("got info.contact.email %v, want me@xeiaso.net", got)
	}

	// An empty security list is how OpenAPI spells "no authentication". A
	// missing key would instead mean "unspecified".
	sec, ok := doc["security"].([]any)
	if !ok {
		t.Fatal("base.json must declare security to state that no auth is required")
	}
	if len(sec) != 0 {
		t.Errorf("got %d security requirements, want 0 (the API takes no auth)", len(sec))
	}

	servers, ok := doc["servers"].([]any)
	if !ok || len(servers) == 0 {
		t.Fatal("base.json must list at least one server")
	}
}

// TestBaseDocumentMerges proves the embedded base satisfies Merge's structural
// expectations, so a bad edit fails here rather than at generate time.
func TestBaseDocumentMerges(t *testing.T) {
	got, err := Merge(baseDoc, []Fragment{
		frag("a.openapi.json", `{"paths": {"/twirp/a.A/Get": {"post": {}}}, "tags": [{"name": "a.A"}]}`),
	})
	if err != nil {
		t.Fatalf("Merge with the embedded base: %v", err)
	}

	doc := map[string]any{}
	if err := json.Unmarshal(got, &doc); err != nil {
		t.Fatalf("merged output is not valid JSON: %v", err)
	}

	if _, ok := doc["paths"].(map[string]any)["/twirp/a.A/Get"]; !ok {
		t.Error("merged document dropped the fragment's path")
	}
}
