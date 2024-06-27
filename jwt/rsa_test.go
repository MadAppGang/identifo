package jwt_test

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	jwti "github.com/madappgang/identifo/v2/jwt"
	"github.com/madappgang/identifo/v2/model"
)

var privateRSAKeyPEM = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC+4Ipbt2/EDJ47
SLDNX97uKxhADjPPvJJlSCCHBf3csU5v9QBxHCUowJf5tMw0R4fR+aVR/ZvZovaB
CkYAyB+wr4inxniWahd0sCVF3+5XvbJ0D1GmizjsNpGfQa6qpXUVn9CdEvOb62a9
psnWx12PM/UBI1FrG7qqyW7TBxnLx2KnsH5kfBW4y+3zonxMbEbBw0bIDvn52gbW
2pbcfogDWxojKIpVJSYgrdT1AISVOYM08RelpcNHmqsHegJ3O+zLHfkTqkfGeIeP
HHMP3BQ1UFgF38bWjh5NzQhyJtUrMb+pu55snzi7UF8zVsfyiAooPOBuCZMEiPkp
52NIeGDPAgMBAAECggEAEuA/rnxAeEHLMA+rNFQjxqfKWSNOal+6lnuAg/nKthVu
rVGsPoNLBXGuVcpUW2Mrgk9O0wHidK5R9Ebgz1j7EUz6laTh7fYF5cs5lGRlvJWM
3T9akr633Vw0IGytakC8iGvqhG4IW0X3PhANa8kBbpTzyK4GcjImzpbm98V+/pDI
Z5NgPeiKY5w7F8o5fGPTBxxSkz0WSGGKFVMbM5JcfrAU3o5ASXxdjbzQPfjdL4n7
/DqqcnxCt2YfzkrRWM/9Mjn0T8CWxI2doqpkQtjjhZ2ob9ybeF3KOdlaWswgffYz
DkzniWfiGt8QqP0nyi5ldRGvaUB9D4vlKg51q7uz+QKBgQDlZdeB589bIdKanW/N
LWGRoSZchQvUzLuikZnRXKUeYgCMCB/5Z98PC2JR0+TngWWnXfMzObhtjQuAyQc+
77wplYFaoKXHqqKSTNYA0yfPw1qnPh4HSuhNE7NbwEEG+6aw9dMHn4meIibP0X9w
8tHXm5hx6V6ngpJhuQxwxwc3ywKBgQDVAyFRcMEZYku8Atccq5IjBOjYvbZ+lf+a
T9787kTt3+YqdhONqKLp7jcDMnG3r0Dv79nhiDh+3J1WY3mIVVIZoIhVizIQHtKX
1IDep8hv9E2JYNFgjIGmvTe54crQv+2cA6koi0zYJivAR/nVB8Mq8m2qZyAk0iZw
BYxkdKsyjQKBgEVdczoDx26uHonEO29WXp9zlC77yCUTt1UkI9fr5L34MmQlfM2k
vA1HivZlVV0vgnaGcSi3Nm5h7O2HXBqK0WHdpFysIRTsIvaMJ1Xeg7ZOQxY5MUlR
PEc6QszmqIMdCz2NR7+RXUKk3wmONrQHqK5CjWk8gPO0BuFn3Dwp4qPbAoGATb8G
uiLdV9Z4rfabbOtyOzXfhrw3j5xP3pKoYMjWf7vo1jaijGGwlJFNou0WdGSS3wA8
FgUSGbuL8av8/7WkcZYWLKLRcvDNDH2TS7ERh0szwaCEyyh9ac5GOKIg1HA42Wi8
pP+y3HGSJmwe05IxucsiG7/oC4hoXxqnU0MB+UECgYA+CrRhLN9/j9Mv6E5+i5B+
Wrw7fvAUPp37yzPxw4uYHt8im0s+M5tVpmPHlION3tTJR5qv3+zyEEgNnivgIx87
3JFlVkawu8M3ZS1USkiP4pauI1NDJ+0TLFIaeRnfjneecBHWPXMGpdDboXYBiXi6
OPcf6cYSsgBxII/5pt4pGw==
-----END PRIVATE KEY-----
`)

var publicRSAKeyPEM = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvuCKW7dvxAyeO0iwzV/e
7isYQA4zz7ySZUgghwX93LFOb/UAcRwlKMCX+bTMNEeH0fmlUf2b2aL2gQpGAMgf
sK+Ip8Z4lmoXdLAlRd/uV72ydA9Rpos47DaRn0GuqqV1FZ/QnRLzm+tmvabJ1sdd
jzP1ASNRaxu6qslu0wcZy8dip7B+ZHwVuMvt86J8TGxGwcNGyA75+doG1tqW3H6I
A1saIyiKVSUmIK3U9QCElTmDNPEXpaXDR5qrB3oCdzvsyx35E6pHxniHjxxzD9wU
NVBYBd/G1o4eTc0IcibVKzG/qbuebJ84u1BfM1bH8ogKKDzgbgmTBIj5KedjSHhg
zwIDAQAB
-----END PUBLIC KEY-----
`)

func TestKeysManualSerializationServiceRSA(t *testing.T) {
	private, err := jwt.ParseRSAPrivateKeyFromPEM(privateRSAKeyPEM)
	if err != nil {
		fmt.Println(string(privateRSAKeyPEM))
		t.Fatalf("Error parsing key from PEM: %v", err)
	}

	public, err := jwt.ParseRSAPublicKeyFromPEM(publicRSAKeyPEM)
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

	b64Pub := []byte(base64.StdEncoding.EncodeToString(genPEM))
	pemString := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s-----END PUBLIC KEY-----\n",
		jwti.Make64ColsString(b64Pub),
	)
	if pemString != string(publicRSAKeyPEM) {
		fmt.Printf("%s\n", pemString)
		t.Fatalf("generated public key PEM and referenced are not equal.")
	}

	genPEMPrivate, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		t.Fatalf("Error creating private PEM: %v", err)
	}
	b64Priv := []byte(base64.StdEncoding.EncodeToString(genPEMPrivate))
	pemStringPrivate := fmt.Sprintf("-----BEGIN PRIVATE KEY-----\n%s-----END PRIVATE KEY-----\n",
		jwti.Make64ColsString(b64Priv),
	)

	if pemStringPrivate != string(privateRSAKeyPEM) {
		fmt.Printf("%s\n", pemStringPrivate)
		t.Fatalf("generated private key PEM and referenced are not equal.")
	}
}

func TestRSAPemMarshalling(t *testing.T) {
	private, alg, err := jwti.LoadPrivateKeyFromPEMString(string(privateRSAKeyPEM))
	if err != nil {
		fmt.Println(string(privateRSAKeyPEM))
		t.Fatalf("Error parsing key from PEM: %v", err)
	}

	if alg != model.TokenSignatureAlgorithmRS256 {
		t.Fatalf("wrong algorithm in PEM: %v", alg)
	}

	result, err := jwti.MarshalPrivateKeyToPEM(private)
	if err != nil {
		t.Fatalf("Error marshaling key to PEM: %v", err)
	}

	if result != string(privateRSAKeyPEM) {
		fmt.Printf("%s\n", result)
		t.Fatalf("generated private key PEM and referenced are not equal.")
	}
}
