package htmx

import (
	"embed"
	"net/http"
)

//go:generate go tool templ generate

var (
	//go:embed *.js
	Static embed.FS
)

func init() {
	Mount(http.DefaultServeMux)
}

// URL is the folder path where the HTMX static files are served from.
const URL = "/.within.website/x/htmx/"

// Mount the HTMX static directory to a given path on a ServeMux.
//
// If you use the Core or Use functions, you will need to ensure that HTMX is mounted at the correct path. Otherwise it will not be able to find the required JavaScript files.
func Mount(mux *http.ServeMux) {
	hdlr := http.StripPrefix(URL, http.FileServer(http.FS(Static)))
	hdlr = UnchangingCache(hdlr)

	mux.Handle(URL, hdlr)
}

// HTTP request headers
const (
	// Request header for the user response to an hx-prompt.
	HeaderPrompt = "HX-Prompt"

	// Request header that is always “true” for HTMX requests.
	HeaderRequest = "Hx-Request"
)

// 286 Stop Polling
//
// HTTP status code that tells HTMX to stop polling from a server response.
//
// For more info, see https://htmx.org/docs/#load_polling
const StatusStopPolling int = 286

// Is returns true if the given request was made by HTMX.
func Is(r *http.Request) bool {
	return r.Header.Get(HeaderRequest) == "true"
}
