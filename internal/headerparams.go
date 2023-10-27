package internal

import (
	"strings"
)

func ParseValueAndParams(value string) map[string]string {
	parts := strings.Split(value, ",")
	vals := make(map[string]string)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		parts := strings.Split(part, ";")
		vals[parts[0]] = strings.Join(parts[1:], ";")
	}

	return vals
}
