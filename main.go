package main

import (
	"fmt"
	"net/http"
	"os"

	"emcoded/emcoded.calendar.be/internal/config"
	"emcoded/emcoded.calendar.be/internal/routes"
	"emcoded/emcoded.calendar.be/internal/server"
)

func main() {
	_, cfgErr := config.LoadAppConfig()
	if cfgErr != nil {
		fmt.Printf("%v\n", cfgErr.Error())
		return
	}

	builder := server.CreateServerBuilder()

	builder.
		Init().
		SetAddress(":8080").
		AddFileServer("/root").
		AddGet("/api/images", routes.GetImages).
		AddGet("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}).
		Build()

	server, serverBuilderErr := builder.Build()
	if serverBuilderErr != nil {
		fmt.Printf("%v\n", serverBuilderErr.Error())
		return
	}

	// If TLS cert and key are provided via env, run TLS. Optionally override
	// the TLS port with TLS_PORT (default 8443).
	cert := os.Getenv("TLS_CERT_FILE")
	key := os.Getenv("TLS_KEY_FILE")
	if cert != "" && key != "" {
		tlsPort := os.Getenv("TLS_PORT")
		if tlsPort == "" {
			tlsPort = "8443"
		}
		server.Addr = ":" + tlsPort
		fmt.Printf("Starting HTTPS server on %s\n", server.Addr)
		if err := server.ListenAndServeTLS(cert, key); err != nil {
			fmt.Printf("server exited: %v\n", err)
		}
		return
	}

	// Fallback to plain HTTP
	fmt.Printf("Starting HTTP server on %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("server exited: %v\n", err)
	}
}
