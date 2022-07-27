package main

import (
	"fmt"
	"os"
	"flag"
	"crypto/tls"
	"context"
	"net/http"
	"os/signal"
	"syscall"
)

var (
	tlscert, tlskey string
)

func getEnv(key , fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func main() {
	
	certfile := getEnv("CERT_FILE", "/etc/certs/cert.pem")
	keyfile := getEnv("KEY_FILE", "/etc/certs/key.pem")

	flag.StringVar(&tlscert, "tlsCertFile", certfile , "File contaains the X509 Certificate for HTTPS")
	flag.StringVar(&tlskey, "tlsKeyFile", keyfile, "The Key file for the HTTPS x509 Certificate")

	flag.Parse()

	certs , err := tls.LoadX509KeyPair(tlscert, tlskey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load Key/Pair certificate\n")
	}

	port := getEnv("PORT", "8080")

	_, err_pn := os.LookupEnv("POD_NAMING")

	if !err_pn {
		fmt.Fprintf(os.Stderr,"no environment variable POD_NAMING provided\n")
		os.Exit(3)	
	}

	server := &http.Server {
		Addr: fmt.Sprintf(":%v", port),
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{certs}},
	}

	gs := myValidServerhandler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", gs.serve)
	server.Handler = mux

	go func() {
		if err := server.ListenAndServeTLS("",""); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start TLS HTTP Service\n: %v", err)
		}
	}()

	debug := getEnv ("DEBUG", "no")
	if debug == "yes" {
		fmt.Fprintf(os.Stdout, "the Server is up and running on port: %s", port)
	}

    signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	fmt.Fprintf(os.Stdout, "Get shutdown signal, sutting down webhook server gracefully...")
	server.Shutdown(context.Background())

}