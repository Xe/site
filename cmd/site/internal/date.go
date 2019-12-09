package internal

import "time"

const iOS13DetriFormat = `2006 M1 2`

// IOS13Detri formats a datestamp like iOS 13 does with the Lojban locale.
func IOS13Detri(t time.Time) string {
	return t.Format(iOS13DetriFormat)
}
