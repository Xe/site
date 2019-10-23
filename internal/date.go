package internal

import "time"

const iOS13DetriFormat = `Y2006 M01 2 Mon`

// IOS13Detri formats a datestamp like iOS 13 does with the Lojban locale.
func IOS13Detri(t time.Time) string {
	return t.Format(iOS13DetriFormat)
}
