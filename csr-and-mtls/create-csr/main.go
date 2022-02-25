package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	subject := pkix.Name{
		CommonName:         "dhakal2021",
		Country:            []string{"US"},
		Province:           []string{"Wyoming"},
		Locality:           []string{"Laramie"},
		Organization:       []string{"University of Wyoming"},
		OrganizationalUnit: []string{"Computer Science"},
	}

	template := x509.CertificateRequest{
		Subject:            subject,
		SignatureAlgorithm: x509.PureEd25519,
		PublicKeyAlgorithm: x509.PublicKeyAlgorithm(x509.PureEd25519),
		EmailAddresses:     []string{"sujan.dhakal019@gmail.com"},
	}

	csrByte, err := x509.CreateCertificateRequest(rand.Reader, &template, privKey)
	if err != nil {
		log.Fatal(err)
	}

	//Encode and save public key.

	pubKF, err := os.Create("ed25519_pub.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer pubKF.Close()

	pxixPubKey, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		log.Fatalln("Error while Marshalling.")
	}

	err = pem.Encode(pubKF, &pem.Block{Type: "PUBLIC KEY", Bytes: pxixPubKey})
	if err != nil {
		log.Fatal(err)
	}

	//Encode and save private key.

	privKF, err := os.Create("ed25519_priv.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer privKF.Close()

	pkcs8PrivateKey, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		log.Fatalln("Error while Marshalling.")
	}

	privateKeyPem := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8PrivateKey,
	}

	err = pem.Encode(privKF, &privateKeyPem)
	if err != nil {
		log.Fatal(err)
	}

	//Encode and save csr.

	csrF, err := os.Create("csr.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer csrF.Close()

	err = pem.Encode(csrF, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrByte})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ed25519 public/private keys pair created and saved in current directory.")
	fmt.Println("CSR created and saved in current directory.")
}
