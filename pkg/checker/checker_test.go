package checker_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestHTTPS(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "Current time: %v", time.Now().Format(time.DateTime))
	})

	p := rand.Intn(49151-1024) + 1024
	addr := fmt.Sprintf("localhost:%d", p)
	server := http.Server{Addr: addr, Handler: r}

	crtFile := "./test_CA/server.crt"
	keyFile := "./test_CA/server.key"

	eChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	go func() {
		eChan <- server.ListenAndServeTLS(crtFile, keyFile)
	}()

	// ca, err := os.ReadFile(pemFile)
	ca, err := os.ReadFile(crtFile)
	require.NoError(t, err)
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(ca)
	require.True(t, ok)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
	}

	t.Run(
		"client with certificate",
		func(t *testing.T) {
			cliwCA := &http.Client{Transport: tr}
			resp, err := cliwCA.Get(fmt.Sprintf("https://%s", addr))
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		},
	)

	t.Run(
		"client without certificate",
		func(t *testing.T) {
			cliwoCA := &http.Client{}
			resp, err := cliwoCA.Get(fmt.Sprintf("https://%s", addr))
			require.Error(t, err)
			require.IsType(t, &url.Error{}, err)
			require.IsType(t, &tls.CertificateVerificationError{}, err.(*url.Error).Err)
			require.Nil(t, resp)
		},
	)

	cancel()
	require.ErrorIs(t, http.ErrServerClosed, <-eChan)
}
