package api

import (
	"encoding/hex"
	"reflect"
	"testing"
)

//TestExtractSignature validages signature extract
//validation text tool
//https://dev-tips.org/Generators/Hash/HMAC
func Test_extractSignature(t *testing.T) {
	d := func(s string) []byte {
		b, _ := hex.DecodeString(s)
		return b
	}
	type args struct {
		b64 string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"non empty string 'test'", args{"SHA-256=Aymga2LNFrM+tnkr6MYLFY2Jou46h2/Omogeu0iMCRQ="}, d("0329A06B62CD16B33EB6792BE8C60B158D89A2EE3A876FCE9A881EBB488C0914")},
		{"empty data", args{""}, nil},
		{"corrupted data", args{"SHA-256=Aymga2LNFrM"}, nil},
		{"no prefix", args{"Aymga2LNFrM+tnkr6MYLFY2Jou46h2/Omogeu0iMCRQ="}, nil},
		{"wrong prefix", args{"SHA-512=Aymga2LNFrM+tnkr6MYLFY2Jou46h2/Omogeu0iMCRQ="}, nil},
		{"empty data", args{"SHA-256="}, nil},
		{"non empty string 'test2'", args{"SHA-256=9TpPNnsorxe8U99HeujuJZCxhfQ51Yz9oD7PBWs/Yjs="}, d("F53A4F367B28AF17BC53DF477AE8EE2590B185F439D58CFDA03ECF056B3F623B")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractSignature(tt.args.b64); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractSignature() = %v, want %v", hex.EncodeToString(got), hex.EncodeToString(tt.want))
			}
		})
	}
}

func Test_validateBodySignature(t *testing.T) {
	d := func(s string) []byte {
		b, _ := hex.DecodeString(s)
		return b
	}

	type args struct {
		body   []byte
		reqMAC []byte
		secret []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"'test' message signature with 'secret'",
			args{[]byte("test"), d("0329A06B62CD16B33EB6792BE8C60B158D89A2EE3A876FCE9A881EBB488C0914"), []byte("secret")},
			false,
		},
		{
			"'test2' message signature with 'secret'",
			args{[]byte("test2"), d("F53A4F367B28AF17BC53DF477AE8EE2590B185F439D58CFDA03ECF056B3F623B"), []byte("secret")},
			false,
		},
		{
			"'test3' message signature with 'secret2'",
			args{[]byte("test3"), d("6530FCF8901618438CACDE3D3C6B660A49DA0BA95601BCDF8D607FC69B68A13D"), []byte("secret2")},
			false,
		},
		{
			"wrong signature",
			args{[]byte("test"), d("6530FCF8901618438CACDE3D3C6B660A49DA0BA95601BCDF8D607FC69B68A13D"), []byte("secret2")},
			true,
		},
		{
			"empty all values signature",
			args{nil, nil, nil},
			true,
		},
		{
			"test with empty secret",
			args{[]byte("test"), d("6530FCF8901618438CACDE3D3C6B660A49DA0BA95601BCDF8D607FC69B68A13D"), nil},
			true,
		},
		{
			"test with empty secret",
			args{[]byte("{ \n\"username\": \"test@madappgang.com\", \"password\": \"secret\", \"scope\": [\"offline\", \"chat\"] }"), d("ee5b46e622d9dccbe939fba26e802c91a56ccb2e03a6bd010479dd0e80243d9d"), []byte("secret")},
			false,
		},
		{
			"empty body with signed with `secret`",
			args{nil, d("F9E66E179B6747AE54108F82F8ADE8B3C25D76FD30AFDE6C395822C530196169"), []byte("secret")},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateBodySignature(tt.args.body, tt.args.reqMAC, tt.args.secret); (err != nil) != tt.wantErr {
				t.Errorf("validateBodySignature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
