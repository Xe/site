package saasproxytoken

import (
	"context"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/protobuf/types/known/emptypb"
	"xeiaso.net/v4/internal/adminpb"
)

type remoteTokenSource struct {
	curr       *oauth2.Token
	lock       sync.Mutex
	remoteURL  string
	httpClient *http.Client
	ptc        adminpb.Patreon
}

func (r *remoteTokenSource) fetchToken() (*oauth2.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := r.ptc.GetToken(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	var tok oauth2.Token
	tok.AccessToken = resp.AccessToken
	tok.TokenType = resp.TokenType
	tok.RefreshToken = resp.RefreshToken
	tok.Expiry = resp.Expiry.AsTime()

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
		ptc:        adminpb.NewPatreonProtobufClient(remoteURL, httpClient),
	}
}
