package internal

import (
	"fmt"
	"testing"
	"time"
)

func TestIOS13Detri(t *testing.T) {
	cases := []struct {
		in  time.Time
		out string
	}{
		{
			in:  time.Date(2019, time.March, 30, 0, 0, 0, 0, time.FixedZone("UTC", 0)),
			out: "2019 M3 30",
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%s -> %s", cs.in.Format(time.RFC3339), cs.out), func(t *testing.T) {
			result := IOS13Detri(cs.in)
			if result != cs.out {
				t.Fatalf("wanted: %s, got: %s", cs.out, result)
			}
		})
	}
}
