package authn

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
)

// Key generation:
//   openssl genrsa 4096 > private.key
//   openssl rsa -pubout < private.key > public.key

func TestJWT(t *testing.T) {
	if err := LoadKey("private.key"); err != nil {
		t.Fatalf("LoadKey() returns err=%+v", err)
	}
	if err := LoadPublicKey("public.key"); err != nil {
		t.Fatalf("LoadPublicKey() returns err=%+v", err)
	}

	token, err := GenerateToken("test")
	if err != nil {
		t.Fatalf("GenerateToken(%s) returns err=%+v", "test", err)
	}
	t.Logf("GenerateToken(%s) returns token=%s", "test", token)

	parsed, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken(%s) returns err=%+v", token, err)
	}

	claims := parsed.Claims.(jwt.MapClaims)
	if _, ok := claims["sub"]; !ok {
		t.Errorf("VerifyToken(%s) returns no sub", token)
	}

	if claims["sub"].(string) != "test" {
		t.Errorf("VerifyToken(%s) returns wrong sub: got=%s, want=%s", token, claims["subject"].(string), "test")
	}
}
