package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"gopkg.in/mxpv/patreon-go.v1"
	"within.website/x/web"
	"xeiaso.net/v4/internal"
	"xeiaso.net/v4/internal/saasproxytoken"
)

var (
	githubSHA             = flag.String("github-sha", "", "GitHub SHA to use for the site")
	miToken               = flag.String("mi-token", "", "Token to use for the mi API")
	patreonSaasProxyURL   = flag.String("patreon-saasproxy-url", "http://patreon-saasproxy/give-token", "URL to use for the patreon saasproxy")
	tailscaleClientID     = flag.String("tailscale-client-id", "", "Tailscale client ID to use")
	tailscaleClientSecret = flag.String("tailscale-client-secret", "", "Tailscale client secret to use")

	regions = []string{"yyz", "bos", "iad", "den", "dfw", "ord", "mia", "phx", "lax", "yul", "gdl", "sjc", "sea", "atl", "ewr", "qro", "ams", "fra", "cdg", "lhr", "mad", "waw", "arn", "gru", "scl", "otp", "eze", "nrt", "bog", "gig", "hkg", "bom", "jnb", "sin", "syd"}
)

func main() {
	flagenv.Parse()
	flag.Parse()
	internal.Slog()

	slog.Info("starting up", "github-sha", *githubSHA)

	/*
		pc, err := NewPatreonClient(hc)
		if err != nil {
			slog.Error("can't create patreon client", "err", err)
		}

		os.MkdirAll("./var", 0700)

		fs, err := lume.New(context.Background(), &lume.Options{
			Branch:        "main",
			Repo:          "https://github.com/Xe/site",
			StaticSiteDir: "lume",
			URL:           "https://xeiaso.net",
			Development:   false,
			PatreonClient: pc,
			DataDir:       "./var",
			MiToken:       *miToken,
		})
		if err != nil {
			log.Fatal(err)
		}

		defer fs.Close()

		var wg sync.WaitGroup

		for _, region := range regions {
			wg.Add(1)
			go func(region string) {
				defer wg.Done()

				if err := uploadSlug(hc, "xedn-"+region, "./var/site.zip"); err != nil {
					slog.Error("error updating", "region", region, "error", err)
				}
			}(region)
		}

		wg.Wait()
	*/
}

func uploadSlug(cli *http.Client, host, fname string) error {
	fin, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fin.Close()

	req, err := http.NewRequest("PUT", "http://"+host+"/xesite/upload", fin)
	if err != nil {
		return err
	}

	slog.Info("uploading", "host", host)

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return web.NewError(http.StatusOK, resp)
	}

	slog.Info("done", "host", host)

	return nil
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
