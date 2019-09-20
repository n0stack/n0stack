package authn

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
)

// Key generation:
//   openssl genrsa 4096 > private.key
//   openssl rsa -pubout < private.key > public.key
func TestJWTWithOpenSSLRSA(t *testing.T) {
	key, err := LoadKey("private.key")
	if err != nil {
		t.Fatalf("LoadKey() returns err=%+v", err)
	}
	pubkey, err := LoadPublicKey("public.key")
	if err != nil {
		t.Fatalf("LoadPublicKey() returns err=%+v", err)
	}
	privateKey = key
	publicKey = pubkey

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

func TestJWTWithSSHECDSA(t *testing.T) {
	key, err := LoadKey("id_ecdsa")
	if err != nil {
		t.Fatalf("LoadKey() returns err=%+v", err)
	}
	pubkey, err := LoadPublicKey("id_ecdsa.pub")
	if err != nil {
		t.Fatalf("LoadPublicKey() returns err=%+v", err)
	}
	privateKey = key
	publicKey = pubkey

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

func TestJWTWithSSHRSA(t *testing.T) {
	key, err := LoadKey("id_rsa")
	if err != nil {
		t.Fatalf("LoadKey() returns err=%+v", err)
	}
	pubkey, err := LoadPublicKey("id_rsa.pub")
	if err != nil {
		t.Fatalf("LoadPublicKey() returns err=%+v", err)
	}
	privateKey = key
	publicKey = pubkey

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
