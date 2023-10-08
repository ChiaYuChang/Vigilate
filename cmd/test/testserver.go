package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/pflag"
)

func main() {
	var port, lifetime int
	var host, crtfile, keyfile string

	pflag.IntVarP(&port, "port", "p", 4000, "port to listen on")
	pflag.StringVarP(&host, "host", "h", "127.0.0.1", "host to listen on")
	pflag.IntVarP(&lifetime, "lifetime", "l", -1,
		"time duration server will keep alive, if lifetime is less than 0, server will close after 300 seconds")
	pflag.StringVarP(&crtfile, "crt", "c", "", "path to certificate file")
	pflag.StringVarP(&keyfile, "key", "k", "", "path to private key file")
	pflag.Parse()

	isHttps := true
	if crtfile == "" || keyfile == "" {
		isHttps = false
	}

	if _, err := os.Stat(crtfile); err != nil {
		log.Fatal("Cannot find certificate file:", crtfile)
		os.Exit(1)
		return
	}

	if _, err := os.Stat(keyfile); err != nil {
		log.Fatal("Cannot find certificate file:", crtfile)
		os.Exit(1)
		return
	}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Current time: %v", time.Now().Format(time.DateTime))
	})

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(lifetime)*time.Second)
	defer cancel()
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	now := time.Now()

	var err error
	if isHttps {
		log.Printf("server listening on %s (HTTPS)\n", server.Addr)
		err = server.ListenAndServeTLS(crtfile, keyfile)
	} else {
		log.Printf("server listening on %s (HTTP)\n", server.Addr)
		err = server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Println("Server stopped w/o errors, uptime: ", time.Since(now))
	os.Exit(0)
}
