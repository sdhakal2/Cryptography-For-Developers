[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-f059dc9a6f8d3a56e377f745f24479a46679e63a5d9fe6f495e02850cd0d8118.svg)](https://classroom.github.com/online_ide?assignment_repo_id=6519411&assignment_repo_type=AssignmentRepo)
# Certificate Authority Homework

For our last homework, you will create the "other half" of our previous assignment - the certificate authority that takes a certificate signing request, validates it, and then sends back a signed certificate.

## Requirements

Your Certificate Authority must meet the following requirements:

1. Accept certificate signing requests through an HTTP API at the `/sign` path.

1. Your HTTP server _must_ use HTTPS, meaning you'll need to generate a _server_ certficate (using the code I've shown you before) to use with the server.

1. Validate that a `Csr-Auth-Code` header is present in the request, and that the value of the header is a value you're expecting. (I'm okay with everyone just checking to see if it's "12345" like it was for the previous homework). If the header isn't present or the code isn't valid, you should return a 401 status code to indicate the request was unauthorized.

1. You _don't_ need to check whether the client sends a `Content-Disposition` header even though my CA did force clients to do that.

1. Check to see whether the body of the request is empty. 

1. Read the bytes contained in the body of the request into a byte slice and then attempt to parse it into an x509.CertificateRequest. You'll have to go from the http body -> PEM bytes -> DER bytes -> x509.CertificateRequest. There is a function in the x509 package to parse the DER bytes into a CertificateRequest struct. If any of these steps fail, return an appropriate Status Code and an error message to the user.

1. Generate a serial number for the new certificate you're about to create - you can use `rand.Int` to get a random number, it should be somewhere between 0 and 2^150.

1. Create a new `x509.Certificate` and set the following fields: `NotBefore`, `NotAfter`, `SerialNumber`, `Subject` (from the CSR), `SignatureAlgorithm` (from the CSR), `EmailAddresses` (from the CSR), `ExtKeyUsage` (set this to a slice of `x509.ExtKeyUsage` with a single element, `x509.ExtKeyUsageClientAuth`), and finally set `IsCA` to `false`.

1. Sign the newly created certificate with the root CA certificate you've generated. You'll need to pass in the _public key_ from the CSR to the `x509.CreateCertificate` function.

1. Send the signed certificate to the client, setting the `Content-Disposition` header to `attachment; filename="signedCertificate.pem"`.


## Generating your Root CA Certificate

I'll put the code that I've been using to generate root CA certificates up on Discord. Feel free to use that to generate your root CA cert that you'll use to sign the CSRs. You can also use the tool to generate a server certificate suitable to use to enable HTTPS on your server.
