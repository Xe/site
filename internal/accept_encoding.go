package internal

import (
	"expvar"
	"net/http"
	"strconv"
	"strings"

	"tailscale.com/metrics"
)

var (
	acceptEncodings = &metrics.LabelMap{Label: "encoding"}

	validEncodings = []string{
		"gzip",
		"x-gzip",
		"deflate",
		"br",
		"identity",
		"snappy",
		"bzip2",
		"lzma",
		"zstd",
	}
)

func init() {
	expvar.Publish("gauge_xesite_accept_encoding", acceptEncodings)
}

func inValidEncodings(enc string) bool {
	for _, validEnc := range validEncodings {
		if enc == validEnc {
			return true
		}
	}
	return false
}

func AcceptEncodingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, enc := range ParseAcceptEncoding(r.Header.Get("Accept-Encoding")) {
			if !inValidEncodings(enc.Encoding) {
				continue
			}
			acceptEncodings.Add(enc.Encoding, 1)
		}

		next.ServeHTTP(w, r)
	})
}

type EncodingQ struct {
	Encoding string
	Q        float64
}

func ParseAcceptEncoding(acptEnc string) []EncodingQ {
	var eqs []EncodingQ

	encQStrs := strings.Split(acptEnc, ",")
	for _, encQStr := range encQStrs {
		trimedEncQStr := strings.Trim(encQStr, " ")

		encQ := strings.Split(trimedEncQStr, ";")
		if len(encQ) == 1 {
			eq := EncodingQ{encQ[0], 1}
			eqs = append(eqs, eq)
		} else {
			qp := strings.Split(encQ[1], "=")
			q, err := strconv.ParseFloat(qp[1], 64)
			if err != nil {
				panic(err)
			}
			eq := EncodingQ{encQ[0], q}
			eqs = append(eqs, eq)
		}
	}
	return eqs
}

type LangQ struct {
	Lang string
	Q    float64
}

func ParseAcceptLanguage(acptLang string) []LangQ {
	var lqs []LangQ

	langQStrs := strings.Split(acptLang, ",")
	for _, langQStr := range langQStrs {
		trimedLangQStr := strings.Trim(langQStr, " ")

		langQ := strings.Split(trimedLangQStr, ";")
		if len(langQ) == 1 {
			lq := LangQ{langQ[0], 1}
			lqs = append(lqs, lq)
		} else {
			qp := strings.Split(langQ[1], "=")
			q, err := strconv.ParseFloat(qp[1], 64)
			if err != nil {
				panic(err)
			}
			lq := LangQ{langQ[0], q}
			lqs = append(lqs, lq)
		}
	}
	return lqs
}
