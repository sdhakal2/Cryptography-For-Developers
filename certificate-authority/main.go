package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"
)

func main() {
	certificate, err := tls.LoadX509KeyPair("server.pem", "serverkey.pem")
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:    ":443",
		Handler: nil,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{certificate},
		},
	}

	http.HandleFunc("/sign", handleSign)

	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}

func handleSign(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(res, "/sign only accept POST request.\n")
	} else {
		csrAuthCode := req.Header.Get("Csr-Auth-Code")
		if csrAuthCode != "12345" {
			res.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(res, "Missing or wrong value for Csr-Auth-Code.\n")
		} else {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(res, "Couldn't read request body\n")
				log.Fatal(err)
			}
			if string(body) == "" {
				res.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(res, "Missing body of the request.\n")
			} else {
				block, _ := pem.Decode(body)
				if block == nil || block.Type != "CERTIFICATE REQUEST" {
					res.WriteHeader(http.StatusUnauthorized)
					fmt.Fprint(res, "Invalid CSR.\n")
				} else {
					certRequest, err := x509.ParseCertificateRequest(block.Bytes)
					if err != nil {
						res.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(res, "Couldn't parse certificate request.\n")
						log.Fatal(err)
					}
					nX := big.NewInt(2)
					nY := big.NewInt(150)
					nM := big.NewInt(0)
					max := big.NewInt(1)
					_ = max.Exp(nX, nY, nM)
					serialNumber, err := rand.Int(rand.Reader, max)
					if err != nil {
						res.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(res, "Couldn't cenerate serial number.\n")
						log.Fatal(err)
					}
					tn := time.Now()
					tf := tn.AddDate(2, 0, 0)
					newCertTemplate := x509.Certificate{
						NotBefore:             tn,
						NotAfter:              tf,
						SerialNumber:          serialNumber,
						Subject:               certRequest.Subject,
						SignatureAlgorithm:    certRequest.SignatureAlgorithm,
						EmailAddresses:        certRequest.EmailAddresses,
						ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
						BasicConstraintsValid: true,
						IsCA:                  false,
					}
					caBuf, err := ioutil.ReadFile("ca.pem")
					if err != nil {
						res.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(res, "Couldn't read root certificate.\n")
						log.Fatal(err)
					}
					caBlock, _ := pem.Decode(caBuf)
					if caBlock == nil || caBlock.Type != "CERTIFICATE" {
						res.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(res, "Invalid or couldn't parse root certificate.\n")
						log.Fatal(err)
					}
					rootCA, err := x509.ParseCertificate(caBlock.Bytes)
					if err != nil {
						res.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(res, "Couldn't parse root certificate.\n")
						log.Fatal(err)
					}
					caKeyBuf, err := ioutil.ReadFile("cakey.pem")
					if err != nil {
						res.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(res, "Couldn't read root certificate key file.\n")
						log.Fatal(err)
					}
					caKeyBlock, _ := pem.Decode(caKeyBuf)
					if caKeyBlock == nil || caKeyBlock.Type != "ED25519 PRIVATE KEY" {
						res.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(res, "Invalid or couldn't parse root certificate key.\n")
						log.Fatal(err)
					}
					caKey, err := x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
					if err != nil {
						res.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(res, "Couldn't parse root certificate key file.\n")
						log.Fatal(err)
					}
					certBuf, err := x509.CreateCertificate(rand.Reader, &newCertTemplate, rootCA, certRequest.PublicKey, caKey)
					if err != nil {
						res.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(res, "Couldn't generate certificate.\n")
						log.Fatal(err)
					}
					res.WriteHeader(http.StatusOK)
					err = pem.Encode(res, &pem.Block{Type: "CERTIFICATE", Bytes: certBuf})
					if err != nil {
						res.WriteHeader(http.StatusInternalServerError)
						fmt.Fprint(res, "Couldn't send certificate.\n")
						log.Fatal(err)
					}
				}
			}
		}
	}
}
