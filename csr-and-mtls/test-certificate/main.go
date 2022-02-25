package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	cert, err := tls.LoadX509KeyPair("../receive-certificate/cert.pem", "../create-csr/ed25519_priv.pem")
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	resp, err := client.Get("https://mtls.invariant.dev")
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode > 299 {
		log.Fatalln("error sending http request, status code:", resp.StatusCode)
	}

	io.Copy(os.Stdout, resp.Body)

}

/*
	NOTE: I mostly copied and past code that Mason provided in dicord.
*/
