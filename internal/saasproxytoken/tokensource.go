package saasproxytoken

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"within.website/x/web"
)

type remoteTokenSource struct {
	curr       *oauth2.Token
	lock       sync.Mutex
	remoteURL  string
	httpClient *http.Client
}

func (r *remoteTokenSource) fetchToken() (*oauth2.Token, error) {
	resp, err := r.httpClient.Get(r.remoteURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, web.NewError(http.StatusOK, resp)
	}

	var tok oauth2.Token
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return nil, err
	}

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
	}
}
