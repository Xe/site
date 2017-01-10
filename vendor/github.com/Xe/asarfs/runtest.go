// +build ignore

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Xe/asarfs"
)

func do404(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
}

func main() {
	fs, err := asarfs.New("./static.asar", http.HandlerFunc(do404))
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":"+os.Getenv("PORT"), fs)
}
