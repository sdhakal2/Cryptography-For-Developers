package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	csrFile, err := os.Open("../create-csr/csr.pem")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	signUrl := "https://crypt.invariant.dev/sign"
	req, err := http.NewRequest("POST", signUrl, csrFile)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Csr-Auth-Code", "12345")
	req.Header.Set("Content-Disposition", "csr.pem")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode > 299 {
		log.Fatalln("Cannot get certificate from the server.")
	}

	certF, err := os.Create("cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer certF.Close()

	io.Copy(certF, resp.Body)

	fmt.Println("Signed certificate reveived and saved in current directory.")
}
