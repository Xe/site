package internal

import "testing"

func TestInValidEncodings(t *testing.T) {
	tests := []struct {
		enc string
		ok  bool
	}{
		{"gzip", true},
		{"x-gzip", true},
		{"tacobell", false},
	}

	for _, test := range tests {
		t.Run(test.enc, func(t *testing.T) {
			ok := inValidEncodings(test.enc)
			if ok != test.ok {
				t.Errorf("ok = %t, want %t", ok, test.ok)
			}
		})
	}
}

func TestParseAcceptLanguage(t *testing.T) {
	acptLang := "en-US,en;q=0.9,ja-JP;q=0.8,ja;q=0.7"
	lqs := ParseAcceptLanguage(acptLang)
	if len(lqs) != 4 {
		t.Errorf("len(lqs) = %d, want 4", len(lqs))
	}
	if lqs[0].Lang != "en-US" {
		t.Errorf("lqs[0].Lang = %s, want en-US", lqs[0].Lang)
	}
	if lqs[0].Q != 1 {
		t.Errorf("lqs[0].Q = %f, want 1", lqs[0].Q)
	}
	if lqs[1].Lang != "en" {
		t.Errorf("lqs[1].Lang = %s, want en", lqs[1].Lang)
	}
	if lqs[1].Q != 0.9 {
		t.Errorf("lqs[1].Q = %f, want 0.9", lqs[1].Q)
	}
	if lqs[2].Lang != "ja-JP" {
		t.Errorf("lqs[2].Lang = %s, want ja-JP", lqs[2].Lang)
	}
	if lqs[2].Q != 0.8 {
		t.Errorf("lqs[2].Q = %f, want 0.8", lqs[2].Q)
	}
	if lqs[3].Lang != "ja" {
		t.Errorf("lqs[3].Lang = %s, want ja", lqs[3].Lang)
	}
	if lqs[3].Q != 0.7 {
		t.Errorf("lqs[3].Q = %f, want 0.7", lqs[3].Q)
	}
}

func TestParseAcceptEncoding(t *testing.T) {
	acptEnc := "gzip, deflate, br;q=1"
	eqs := ParseAcceptEncoding(acptEnc)
	if len(eqs) != 3 {
		t.Errorf("len(eqs) = %d, want 3", len(eqs))
	}
	if eqs[0].Encoding != "gzip" {
		t.Errorf("eqs[0].Encoding = %s, want gzip", eqs[0].Encoding)
	}
	if eqs[0].Q != 1 {
		t.Errorf("eqs[0].Q = %f, want 1", eqs[0].Q)
	}
	if eqs[1].Encoding != "deflate" {
		t.Errorf("eqs[1].Encoding = %s, want deflate", eqs[1].Encoding)
	}
	if eqs[1].Q != 1 {
		t.Errorf("eqs[1].Q = %f, want 1", eqs[1].Q)
	}
	if eqs[2].Encoding != "br" {
		t.Errorf("eqs[2].Encoding = %s, want br", eqs[2].Encoding)
	}
	if eqs[2].Q != 1 {
		t.Errorf("eqs[2].Q = %f, want 1", eqs[2].Q)
	}
}
