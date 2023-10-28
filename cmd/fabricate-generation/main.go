package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/facebookgo/flagenv"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"gopkg.in/mxpv/patreon-go.v1"
	"tailscale.com/client/tailscale"
	"tailscale.com/hostinfo"
	"tailscale.com/ipn/store/mem"
	"tailscale.com/tsnet"
	"tailscale.com/util/cmpx"
	"within.website/x/web"
	"xeiaso.net/v4/internal"
	"xeiaso.net/v4/internal/lume"
	"xeiaso.net/v4/internal/saasproxytoken"
)

var (
	githubSHA             = flag.String("github-sha", "", "GitHub SHA to use for the site")
	miToken               = flag.String("mi-token", "", "Token to use for the mi API")
	patreonSaasProxyURL   = flag.String("patreon-saasproxy-url", "http://patreon-saasproxy/give-token", "URL to use for the patreon saasproxy")
	tailscaleClientID     = flag.String("tailscale-client-id", "", "Tailscale client ID to use")
	tailscaleClientSecret = flag.String("tailscale-client-secret", "", "Tailscale client secret to use")
)

func main() {
	// Required to use the Tailscale client API. This is sussy, but okay.
	tailscale.I_Acknowledge_This_API_Is_Unstable = true

	flagenv.Parse()
	flag.Parse()
	internal.Slog()

	baseURL := cmpx.Or(os.Getenv("TS_BASE_URL"), "https://api.tailscale.com")

	credentials := clientcredentials.Config{
		ClientID:     *tailscaleClientID,
		ClientSecret: *tailscaleClientSecret,
		TokenURL:     baseURL + "/api/v2/oauth/token",
		Scopes:       []string{"device"},
	}

	ctx := context.Background()
	tsClient := tailscale.NewClient("-", nil)
	tsClient.HTTPClient = credentials.Client(ctx)
	tsClient.BaseURL = baseURL

	caps := tailscale.KeyCapabilities{
		Devices: tailscale.KeyDeviceCapabilities{
			Create: tailscale.KeyDeviceCreateCapabilities{
				Reusable:      false,
				Ephemeral:     true,
				Preauthorized: true,
				Tags:          []string{"tag:service", "tag:ci"},
			},
		},
	}

	authkey, _, err := tsClient.CreateKey(ctx, caps)
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Args[0] = "via XeDN"

	hostinfo.SetApp("xeiaso.net/v4/cmd/fabricate-generation")

	memStore, err := mem.New(log.Printf, "")
	if err != nil {
		log.Fatal(err)
	}

	srv := &tsnet.Server{
		Hostname:  "github-action-" + (*githubSHA)[:7],
		Logf:      log.Printf,
		Ephemeral: true,
		Store:     memStore,
		AuthKey:   authkey,
	}

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	if _, err := srv.Up(context.Background()); err != nil {
		log.Fatal(err)
	}

	hc := srv.HTTPClient()

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

	for _, region := range []string{"fra", "sea", "yyz"} {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()

			if err := uploadSlug(hc, "xedn-"+region, "./var/site.zip"); err != nil {
				slog.Error("error updating", "region", region, "error", err)
			}
		}(region)
	}

	wg.Wait()
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
