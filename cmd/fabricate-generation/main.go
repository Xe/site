package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"gopkg.in/mxpv/patreon-go.v1"
	"within.website/x/tigris"
	"xeiaso.net/v4/internal"
	"xeiaso.net/v4/internal/lume"
	"xeiaso.net/v4/internal/saasproxytoken"
)

var (
	bucketName          = flag.String("bucket-name", "xesite", "Name of the S3 bucket to upload to")
	githubSHA           = flag.String("github-sha", "", "GitHub SHA to use for the site")
	miURL               = flag.String("mimi-announce-url", "", "Mi url")
	patreonSaasProxyURL = flag.String("patreon-saasproxy-url", "http://xesite-patreon-saasproxy.flycast", "URL to use for the patreon saasproxy")
	siteURL             = flag.String("site-url", "https://xeiaso.net/", "URL to use for the site")
)

func main() {
	flagenv.Parse()
	flag.Parse()
	internal.Slog()

	slog.Info("starting up", "github-sha", *githubSHA)

	pc, err := NewPatreonClient(http.DefaultClient)
	if err != nil {
		slog.Error("can't create patreon client", "err", err)
	}

	os.MkdirAll("./var", 0700)

	s3c, err := tigris.Client(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	_ = s3c

	fs, err := lume.New(context.Background(), &lume.Options{
		Branch:        "main",
		Repo:          "https://github.com/Xe/site",
		StaticSiteDir: "lume",
		URL:           *siteURL,
		Development:   false,
		PatreonClient: pc,
		DataDir:       "./var",
		MiURL:         *miURL,
	})
	if err != nil {
		log.Fatal(err)
	}

	defer fs.Close()

	// if err := uploadFolderToS3(context.Background(), s3c, "./var/repo/lume/_site", *bucketName); err != nil {
	// 	log.Fatal(err)
	// }
}

func NewPatreonClient(hc *http.Client) (*patreon.Client, error) {
	ts := saasproxytoken.RemoteTokenSource(*patreonSaasProxyURL, hc)
	tc := oauth2.NewClient(context.Background(), ts)

	client := patreon.NewClient(tc)
	if u, err := client.FetchUser(); err != nil {
		return nil, err
	} else {
		slog.Info("logged in as", "user", u.Data.Attributes.FullName)
	}

	return client, nil
}

func uploadFolderToS3(ctx context.Context, s3c *s3.Client, folderPath, bucketName string) error {
	// Ensure folderPath ends with a slash to correctly trim the prefix from file paths
	cleanFolderPath := filepath.Clean(folderPath) + string(os.PathSeparator)

	err := filepath.Walk(cleanFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file %q, %v", path, err)
			}
			defer file.Close()

			fileContent, err := io.ReadAll(file)
			if err != nil {
				return fmt.Errorf("failed to read file %q, %v", path, err)
			}

			ext := filepath.Ext(path)
			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				// otherwise, detect by content
				mimeType = http.DetectContentType(fileContent)
			}

			key := strings.TrimPrefix(path, cleanFolderPath)

			_, err = s3c.PutObject(ctx, &s3.PutObjectInput{
				Bucket:      aws.String(bucketName),
				Key:         aws.String(key),
				Body:        bytes.NewReader(fileContent),
				ContentType: aws.String(mimeType),
			})
			if err != nil {
				return fmt.Errorf("failed to upload file %q to bucket %q, %v", path, bucketName, err)
			}
			slog.Info("uploaded file", "file", path, "bucket", bucketName)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk through folder %q, %v", folderPath, err)
	}

	return nil
}
