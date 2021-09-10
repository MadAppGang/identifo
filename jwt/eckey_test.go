package jwt_test

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"testing"

	jwt "github.com/golang-jwt/jwt/v4"
	jwti "github.com/madappgang/identifo/jwt"
	"github.com/madappgang/identifo/model"
)

// SEC 1, ASN.1, DER format
// 	privateKeyPEM := []byte(`-----BEGIN EC PRIVATE KEY-----
// MHcCAQEEIKcw4Osfw4a5G11fWprAjxrPLSAhKv5H5Gj27NBXsGDKoAoGCCqGSM49
// AwEHoUQDQgAED3DoOWZMbqYc0OO1Ih628hB2Odhv4mjl1vt0iBu3gTKz1XAk+YHG
// 8aoI+42TJle6hawmTzrSD6khaNRaDQAKbg==
// -----END EC PRIVATE KEY-----
// `)

var privateKeyPEM = []byte(`-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgpzDg6x/DhrkbXV9a
msCPGs8tICEq/kfkaPbs0FewYMqhRANCAAQPcOg5ZkxuphzQ47UiHrbyEHY52G/i
aOXW+3SIG7eBMrPVcCT5gcbxqgj7jZMmV7qFrCZPOtIPqSFo1FoNAApu
-----END PRIVATE KEY-----
`)

var publicKeyPEM = []byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAED3DoOWZMbqYc0OO1Ih628hB2Odhv
4mjl1vt0iBu3gTKz1XAk+YHG8aoI+42TJle6hawmTzrSD6khaNRaDQAKbg==
-----END PUBLIC KEY-----
`)

func TestKeysManualSerializationService(t *testing.T) {
	private, err := jwt.ParseECPrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		fmt.Println(string(privateKeyPEM))
		t.Fatalf("Error parsing key from PEM: %v", err)
	}

	public, err := jwt.ParseECPublicKeyFromPEM(publicKeyPEM)
	if err != nil {
		t.Fatalf("Error parsing key from PEM: %v", err)
	}

	generatedPublic := private.Public()

	if !public.Equal(generatedPublic) {
		t.Fatalf("Generated public key and reference keys are not equal")
	}

	genPEM, err := x509.MarshalPKIXPublicKey(generatedPublic)
	if err != nil {
		t.Fatalf("Error creating PEM: %v", err)
	}
	// creating the same certificate as ssl tool is doing with the following command:
	// openssl ec -in private.pem -pubout -out public.pem
	pemString := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n%s\n-----END PUBLIC KEY-----\n",
		base64.StdEncoding.EncodeToString(genPEM)[0:64],
		base64.StdEncoding.EncodeToString(genPEM)[64:],
	)
	if pemString != string(publicKeyPEM) {
		fmt.Printf("%s\n", pemString)
		t.Fatalf("generated public key PEM and referenced are not equal.")
	}

	genPEMPrivate, err := x509.MarshalPKCS8PrivateKey(private)
	//	genPEMPrivate, err := x509.MarshalECPrivateKey(private)
	if err != nil {
		t.Fatalf("Error creating private PEM: %v", err)
	}
	pemStringPrivate := fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s\n%s\n%s\n-----END PRIVATE KEY-----\n",
		base64.StdEncoding.EncodeToString(genPEMPrivate)[0:64],
		base64.StdEncoding.EncodeToString(genPEMPrivate)[64:128],
		base64.StdEncoding.EncodeToString(genPEMPrivate)[128:],
	)

	if pemStringPrivate != string(privateKeyPEM) {
		fmt.Printf("%s\n", pemStringPrivate)
		t.Fatalf("generated private key PEM and referenced are not equal.")
	}
}

func TestPemMarshalling(t *testing.T) {
	// private, err := jwt.ParseECPrivateKeyFromPEM(privateKeyPEM)
	private, alg, err := jwti.LoadPrivateKeyFromString(string(privateKeyPEM))
	if err != nil {
		fmt.Println(string(privateKeyPEM))
		t.Fatalf("Error parsing key from PEM: %v", err)
	}

	if alg != model.TokenSignatureAlgorithmES256 {
		t.Fatalf("wrong algorithm in PEM: %v", alg)
	}

	result, err := jwti.MarshalPrivateKeyToPEM(private)
	if err != nil {
		t.Fatalf("Error marshaling key to PEM: %v", err)
	}

	if result != string(privateKeyPEM) {
		fmt.Printf("%s\n", result)
		t.Fatalf("generated private key PEM and referenced are not equal.")
	}
}
