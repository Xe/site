// +build go1.8

package asarfs

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"testing"
)

func BenchmarkHTTPFileSystem(b *testing.B) {
	fs := http.FileServer(http.Dir("."))

	l, s, err := setupHandler(fs)
	if err != nil {
		b.Fatal(err)
	}
	defer l.Close()
	defer s.Close()

	url := fmt.Sprintf("http://%s", l.Addr())

	for n := 0; n < b.N; n++ {
		testHandler(url)
	}
}

func BenchmarkASARfs(b *testing.B) {
	fs, err := New("./static.asar", http.HandlerFunc(do404))
	if err != nil {
		b.Fatal(err)
	}

	l, s, err := setupHandler(fs)
	if err != nil {
		b.Fatal(err)
	}
	defer l.Close()
	defer s.Close()

	url := fmt.Sprintf("http://%s", l.Addr())

	for n := 0; n < b.N; n++ {
		testHandler(url)
	}
}

func BenchmarkPreloadedASARfs(b *testing.B) {
	for n := 0; n < b.N; n++ {
		testHandler(asarfsurl)
	}
}

func BenchmarkASARfsHTTPFilesystem(b *testing.B) {
	fs, err := New("./static.asar", http.HandlerFunc(do404))
	if err != nil {
		b.Fatal(err)
	}

	l, s, err := setupHandler(http.FileServer(fs))
	if err != nil {
		b.Fatal(err)
	}
	defer l.Close()
	defer s.Close()

	url := fmt.Sprintf("http://%s", l.Addr())

	for n := 0; n < b.N; n++ {
		testHandler(url)
	}
}

func BenchmarkPreloadedASARfsHTTPFilesystem(b *testing.B) {
	for n := 0; n < b.N; n++ {
		testHandler(asarfshttpfsurl)
	}
}

func do404(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
}

func setupHandler(h http.Handler) (net.Listener, *http.Server, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	s := &http.Server{
		Handler: h,
	}
	go s.ListenAndServe()

	return l, s, nil
}

func testHandler(u string) error {
	num := rand.Intn(9)
	num++
	sub := rand.Intn(99)

	fname := fmt.Sprintf("/static/%d/%d%d.json", num, num, sub)

	resp, err := http.Get(u + fname)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		panic(err)
	}

	return nil
}

var (
	asarfsurl       string
	asarfshttpfsurl string
)

func TestMain(m *testing.M) {
	go func() {
		fs, err := New("./static.asar", http.HandlerFunc(do404))
		if err != nil {
		}

		l, _, err := setupHandler(fs)
		if err != nil {
		}

		asarfsurl = fmt.Sprintf("http://%s", l.Addr().String())
	}()

	go func() {
		fs, err := New("./static.asar", http.HandlerFunc(do404))
		if err != nil {
		}

		l, _, err := setupHandler(http.FileServer(fs))
		if err != nil {
		}

		asarfshttpfsurl = fmt.Sprintf("http://%s", l.Addr().String())
	}()

	os.Exit(m.Run())
}
