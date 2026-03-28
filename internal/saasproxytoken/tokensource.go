package saasproxytoken

import (
	"context"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	adminpb "xeiaso.net/v4/gen/xeiaso/net/admin/v1"
)

type remoteTokenSource struct {
	curr       *oauth2.Token
	lock       sync.Mutex
	remoteURL  string
	httpClient *http.Client
	ptc        adminpb.PatreonService
}

func (r *remoteTokenSource) fetchToken() (*oauth2.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := r.ptc.GetToken(ctx, &adminpb.GetTokenRequest{})
	if err != nil {
		return nil, err
	}

	var tok oauth2.Token
	tok.AccessToken = resp.Token.AccessToken
	tok.TokenType = resp.Token.TokenType
	tok.RefreshToken = resp.Token.RefreshToken
	tok.Expiry = resp.Token.Expiry.AsTime()

	return &tok, nil
}

func (r *remoteTokenSource) Token() (*oauth2.Token, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.curr == nil {
		tok, err := r.fetchToken()
		if err != nil {
			return nil, err
		}
		r.curr = tok
		return tok, nil
	}

	if r.curr.Expiry.Before(time.Now()) {
		tok, err := r.fetchToken()
		if err != nil {
			return nil, err
		}
		r.curr = tok
	}

	return r.curr, nil
}

func RemoteTokenSource(remoteURL string, httpClient *http.Client) oauth2.TokenSource {
	return &remoteTokenSource{
		remoteURL:  remoteURL,
		httpClient: httpClient,
		ptc:        adminpb.NewPatreonServiceProtobufClient(remoteURL, httpClient),
	}
}
