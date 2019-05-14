package jwt

import (
	"reflect"
	"testing"

	"github.com/madappgang/identifo/mem"
	"github.com/madappgang/identifo/model"
)

const (
	privateKey         = "./private.pem"
	publicKey          = "./public.pem"
	testIssuer         = "identifo.madappgang.com"
	tokenStringExample = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsInN1YiI6IjEyMzQ1Njc4OTAifQ.Sqmh_44nXg3Lxs9jr9YCDZVNJN459Br4ODnZIt3EY72opwy5hzYL_l_hua4PJCM0WmYNLB-nKC80TS84LO5muw"
)

func TestNewTokenService(t *testing.T) {
	us, err := mem.NewUserStorage()
	if err != nil {
		t.Errorf("Unable to create user storage %v", err)
	}
	tstor, err := mem.NewTokenStorage()
	if err != nil {
		t.Errorf("Unable to create token storage %v", err)
	}
	as, err := mem.NewAppStorage()
	if err != nil {
		t.Errorf("Unable to create app storage %v", err)
	}
	ts, err := NewTokenService(privateKey, publicKey, testIssuer, model.TokenServiceAlgorithmES256, tstor, as, us)
	if err != nil {
		t.Errorf("Unable to crate service %v", err)
	}
	type args struct {
		private string
		public  string
		issuer  string
	}
	tests := []struct {
		name    string
		args    args
		want    model.TokenService
		wantErr bool
	}{
		{"successfull creation", args{privateKey, publicKey, testIssuer}, ts, false},
		{"invalid private path", args{"somepath", publicKey, testIssuer}, nil, true},
		{"invalid public path", args{privateKey, "domefakepath", testIssuer}, nil, true},
		{"empty file pathes", args{"", "", testIssuer}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTokenService(tt.args.private, tt.args.public, testIssuer, model.TokenServiceAlgorithmES256, tstor, as, us)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTokenService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseString(t *testing.T) {
	us, err := mem.NewUserStorage()
	if err != nil {
		t.Errorf("Unable to create user storage %v", err)
	}
	tstor, err := mem.NewTokenStorage()
	if err != nil {
		t.Errorf("Unable to create token storage %v", err)
	}
	as, err := mem.NewAppStorage()
	if err != nil {
		t.Errorf("Unable to create app storage %v", err)
	}
	ts, err := NewTokenService(privateKey, publicKey, testIssuer, model.TokenServiceAlgorithmES256, tstor, as, us)
	if err != nil {
		t.Errorf("Unable to crate service %v", err)
	}
	token, err := ts.Parse(tokenStringExample)
	if err != nil {
		t.Errorf("Unable to parse token %v", err)
	}
	if token == nil {
		t.Error("Token is empty")
	}

	tkn, ok := token.(*Token)
	if !ok {
		t.Error("Token is wrong type")
	}
	claims, ok := tkn.JWT.Claims.(*Claims)
	if !ok {
		t.Error("Claims are invalid")
	}
	if claims.Subject != "1234567890" {
		t.Errorf("Claims subject is invalid, got %v, want: %v", claims.Subject, "1234567890")
	}
	if claims.IssuedAt != 1516239022 {
		t.Errorf("Claims issued At is invalid, got %v, want: %v", claims.IssuedAt, 1516239022)
	}

}

func TestTokenToString(t *testing.T) {
	us, err := mem.NewUserStorage()
	if err != nil {
		t.Errorf("Unable to create user storage %v", err)
	}
	tstor, err := mem.NewTokenStorage()
	if err != nil {
		t.Errorf("Unable to create token storage %v", err)
	}
	as, err := mem.NewAppStorage()
	if err != nil {
		t.Errorf("Unable to create app storage %v", err)
	}
	ts, err := NewTokenService(privateKey, publicKey, testIssuer, model.TokenServiceAlgorithmES256, tstor, as, us)
	if err != nil {
		t.Errorf("Unable to create service %v", err)
	}
	token, err := ts.Parse(tokenStringExample)
	if err != nil {
		t.Errorf("Unable to parse token %v", err)
	}
	if token == nil {
		t.Error("Token is empty")
	}

	tokenString, err := ts.String(token)
	if err != nil {
		t.Errorf("Unable to serialize token %v", err)
	}
	if tokenString == tokenStringExample {
		t.Errorf("Generated token is matched, should not, generated: %v, expected: %v", tokenString, tokenStringExample)
	}
	token2, err := ts.Parse(tokenString)
	if err != nil {
		t.Errorf("Unable to parse token %v", err)
	}
	if token2 == nil {
		t.Error("Token is empty")
	}
	t1, _ := token.(*Token)
	t2, _ := token2.(*Token)
	claims1, _ := t1.JWT.Claims.(*Claims)
	claims2, _ := t2.JWT.Claims.(*Claims)

	if !reflect.DeepEqual(t1.JWT.Header, t2.JWT.Header) {
		t.Errorf("Headers = %+v, want %+v", t1.JWT.Header, t2.JWT.Header)
	}
	if !reflect.DeepEqual(claims1, claims2) {
		t.Errorf("Claims = %+v, want %+v", claims1, claims2)
	}

}

func TestNewToken(t *testing.T) {
	us, err := mem.NewUserStorage()
	if err != nil {
		t.Errorf("Unable to create user storage %v", err)
	}
	tstor, err := mem.NewTokenStorage()
	if err != nil {
		t.Errorf("Unable to create token storage %v", err)
	}
	as, err := mem.NewAppStorage()
	if err != nil {
		t.Errorf("Unable to create app storage %v", err)
	}
	ts, err := NewTokenService(privateKey, publicKey, testIssuer, model.TokenServiceAlgorithmES256, tstor, as, us)
	if err != nil {
		t.Errorf("Unable to create service %v", err)
	}
	ustg, _ := mem.NewUserStorage()
	user, _ := ustg.UserByNamePassword("username", "password")
	//generate random user until we get active user
	for !user.Active() {
		user, _ = ustg.UserByNamePassword("username", "password")
	}
	scopes := []string{"scope1", "scope2"}
	tokenPayload := []string{"name"}
	app := mem.MakeAppData("123456", "1", true, "testName", "testDescriprion", scopes, true, "", 0, 0, 0, tokenPayload, true)
	token, err := ts.NewToken(user, scopes, app)
	if err != nil {
		t.Errorf("Unable to create token %v", err)
	}
	tokenString, err := ts.String(token)
	if err != nil {
		t.Errorf("Unable to serialize token %v", err)
	}
	token2, err := ts.Parse(tokenString)
	if err != nil {
		t.Errorf("Unable to parse token %v", err)
	}
	if token2 == nil {
		t.Error("Token is empty")
	}
	t2, _ := token2.(*Token)
	claims2, _ := t2.JWT.Claims.(*Claims)
	if _, ok := claims2.Payload["name"]; !ok {
		t.Errorf("Claims payload = %+v, want name in payload.", claims2.Payload)
	}
	if claims2.Issuer != testIssuer {
		t.Errorf("Issuer = %+v, want %+v", claims2.Issuer, testIssuer)
	}
	if claims2.Subject != user.ID() {
		t.Errorf("Subject = %+v, want %+v", claims2.Subject, user.ID())
	}
	if claims2.Audience != app.ID() {
		t.Errorf("Audience = %+v, want %+v", claims2.Audience, app.ID())
	}

}
