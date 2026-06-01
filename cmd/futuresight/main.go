// Command futuresight is the xesite preview server. It stores each build's erofs
// volume in Tigris and serves every version on its own subdomain: a permanent
// host keyed by the truncated sha256 of the volume, and a moving host keyed by a
// DNS-safe slug of the git branch.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"expvar"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	storage "github.com/tigrisdata/storage-go"
	"xeiaso.net/v4/internal"
)

// hashLen must match internal/lume's hashLen: the number of hex characters that
// content-address a volume (and form a permanent subdomain label).
const hashLen = 16

var (
	bind            = flag.String("bind", ":3000", "Port to listen on")
	internalAPIBind = flag.String("internal-api-bind", ":3001", "Port to listen on for the internal API")
	baseDomain      = flag.String("base-domain", "preview.xeiaso.net", "Base domain; the leftmost label selects the preview")
	cacheDir        = flag.String("cache-dir", "./var/volumes", "Directory to cache downloaded erofs volumes")
	tigrisBucket    = flag.String("tigris-bucket", "", "Tigris bucket to store preview volumes in")
	accessKeyID     = flag.String("tigris-access-key-id", "", "Tigris access key ID (falls back to AWS_ACCESS_KEY_ID)")
	secretAccessKey = flag.String("tigris-secret-access-key", "", "Tigris secret access key (falls back to AWS_SECRET_ACCESS_KEY)")
	uploadToken     = flag.String("upload-token", "", "Bearer token required to POST /upload")
	maxUploadBytes  = flag.Int64("max-upload-bytes", 512<<20, "Maximum accepted upload size in bytes")
	generateToken   = flag.Bool("generate-token", false, "Print a random upload token and exit")
)

// randomToken returns a 32-byte cryptographically random token as hex.
func randomToken() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("futuresight: can't generate token: %w", err)
	}
	return hex.EncodeToString(b[:]), nil
}

func main() {
	flagenv.Parse()
	flag.Parse()
	internal.Slog()

	if *generateToken {
		token, err := randomToken()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(token)
		return
	}

	ctx := context.Background()

	if *tigrisBucket == "" {
		log.Fatal("--tigris-bucket is required")
	}

	opts := []storage.Option{storage.WithGlobalEndpoint()}
	if *accessKeyID != "" && *secretAccessKey != "" {
		opts = append(opts, storage.WithAccessKeypair(*accessKeyID, *secretAccessKey))
	}

	client, err := storage.New(ctx, opts...)
	if err != nil {
		log.Fatal(err)
	}

	store, err := NewStore(client, *tigrisBucket, *cacheDir)
	if err != nil {
		log.Fatal(err)
	}

	srv := &server{
		store:          store,
		baseDomain:     *baseDomain,
		uploadToken:    *uploadToken,
		maxUploadBytes: *maxUploadBytes,
	}

	go internalAPI()

	mux := http.NewServeMux()
	mux.HandleFunc("/upload", srv.handleUpload)
	mux.HandleFunc("/", srv.handleServe)

	ln, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("starting futuresight", "bind", *bind, "base-domain", *baseDomain, "bucket", *tigrisBucket)
	log.Fatal(http.Serve(ln, mux))
}

func internalAPI() {
	mux := http.NewServeMux()
	mux.Handle("/debug/vars", expvar.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	ln, err := net.Listen("tcp", *internalAPIBind)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.Serve(ln, mux); err != nil {
		log.Fatal(err)
	}
}
