package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/acme/autocert"
)

var (
	HTML_FILE     = flag.String("html", "", "the html bot file")
	HTTP_SERVER   = flag.String("http", ":http", "the http listen address")
	HTTPS_SERVER  = flag.String("https", ":https", "the https listen address")
	SSL_CACHE_DIR = flag.String("cache-dir", ".autocert", "the autocert cache directory")
	SERVER_NAME   = flag.String("server-name", "", "")
)

func main() {
	// parse the command line flags
	flag.Parse()

	// open the specified html file to be parsed
	file, err := os.Open(*HTML_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Initialize the parser
	bot, err := NewBotFromReader(file)
	if err != nil {
		log.Fatal(err)
	}

	// add a messenger handler
	http.Handle("/messenger", NewMessenger(bot))
	http.Handle("/", func(res http.ResponseWriter, req *http.Request){
		res.Write([]byte("Hello, Bkit ;)"))
	})

	// our global error channel
	errchan := make(chan error)

	// starts the HTTPS server if required
	if *HTTPS_SERVER != "" {
		m := &autocert.Manager{
			Cache:      autocert.DirCache(*SSL_CACHE_DIR),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(strings.Split(*SERVER_NAME, ",")...),
		}
		s := &http.Server{
			Addr:      *HTTPS_SERVER,
			TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		}
		go (func() {
			errchan <- s.ListenAndServeTLS("", "")
		})()
	}

	// starts the HTTP server if required
	if *HTTP_SERVER != "" {
		go (func() {
			errchan <- http.ListenAndServe(*HTTP_SERVER, m.HTTPHandler(nil))
		})()
	}

	// panic with the errors
	log.Fatal(<-errchan)
}
