package k8s

import (
	"context"
	"encoding/json"
	"errors"
	"expvar"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const secretValue = "oauth2-token"

var (
	tokenRefreshCount = expvar.NewInt("gauge_xesite_k8s_token_refresh_count")

	ErrSecretValueDoesntExist = errors.New("internal: can't find oauth2-token in secret")
)

type tokenSource struct {
	secretName string
	namespace  string
	base       oauth2.TokenSource
	clientSet  kubernetes.Interface
}

func (kts *tokenSource) Token() (tok *oauth2.Token, err error) {
	tok, _ = loadTokenFromK8s(context.Background(), kts.clientSet, kts.namespace, kts.secretName)

	if tok != nil && tok.Expiry.After(time.Now()) {
		return tok, nil
	}

	if tok, err = kts.base.Token(); err != nil {
		return nil, err
	}

	tokenRefreshCount.Add(1)

	if err := kts.saveToken(tok); err != nil {
		return nil, err
	}

	return tok, err
}

func (kts *tokenSource) saveToken(tok *oauth2.Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sec, err := kts.clientSet.CoreV1().Secrets(kts.namespace).Get(ctx, kts.secretName, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("can't get secret %s::%s: %w", kts.namespace, kts.secretName, err)
	}

	data, err := json.Marshal(tok)
	if err != nil {
		return fmt.Errorf("can't marshal json: %w", err)
	}

	sec.Data[secretValue] = data

	if _, err := kts.clientSet.CoreV1().Secrets(kts.namespace).Update(ctx, sec, v1.UpdateOptions{}); err != nil {
		return fmt.Errorf("can't update secret %s::%s: %w", kts.namespace, kts.secretName, err)
	}

	return nil
}

func loadTokenFromK8s(ctx context.Context, cs kubernetes.Interface, ns, secretName string) (*oauth2.Token, error) {
	sec, err := cs.CoreV1().Secrets(ns).Get(ctx, secretName, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("can't get secret %s::%s: %w", ns, secretName, err)
	}

	data, ok := sec.Data[secretValue]
	if !ok {
		return nil, ErrSecretValueDoesntExist
	}

	var tok oauth2.Token
	if err := json.Unmarshal(data, &tok); err != nil {
		return nil, fmt.Errorf("can't unmarshal oauth2-token in %s::%s: %w", ns, secretName, err)
	}

	return &tok, nil
}

func TokenSource(namespace string, secretName string, config *oauth2.Config) (oauth2.TokenSource, error) {
	cs, err := getClientSet()
	if err != nil {
		return nil, fmt.Errorf("can't get client set: %w", err)
	}

	tok, err := loadTokenFromK8s(context.Background(), cs, namespace, secretName)
	if err != nil {
		return nil, err
	}

	orig := config.TokenSource(context.Background(), tok)

	return oauth2.ReuseTokenSource(nil, &tokenSource{
		secretName: secretName,
		namespace:  namespace,
		base:       orig,
		clientSet:  cs,
	}), nil
}

func getClientSet() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
		if err != nil {
			return nil, err
		}
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientSet, nil
}
