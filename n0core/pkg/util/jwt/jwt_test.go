package jwt

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
)

// Key generation:
//   openssl genrsa 4096 > private.key
//   openssl rsa -pubout < private.key > public.key
func TestJWTWithOpenSSLRSA(t *testing.T) {
	testAuthentication(t, "private.key", "public.key")
}

// ssh-keygen -t ecdsa
func TestJWTWithSSHECDSA(t *testing.T) {
	testAuthentication(t, "id_ecdsa", "id_ecdsa.pub")
}

// ssh-keygen -t rsa
func TestJWTWithSSHRSA(t *testing.T) {
	testAuthentication(t, "id_rsa", "id_rsa.pub")
}

func testAuthentication(t *testing.T, keyFile, pubkeyFile string) {
	key, err := ParsePrivateKeyFromFile(keyFile)
	if err != nil {
		t.Fatalf("ParsePrivateKeyFromFile() returns err=%+v", err)
	}
	pubkey, err := ParsePublicKeyFromFile(pubkeyFile)
	if err != nil {
		t.Fatalf("ParsePublicKeyFromFile() returns err=%+v", err)
	}

	chaltoken, cookie, err := key.GenerateChallengeToken("test")
	if err != nil {
		t.Fatalf("GenerateAuthenticationToken(%s) returns err=%+v", "test", err)
	}
	t.Logf("GenerateAuthenticationToken(%s) returns token=%s", "test", chaltoken)

	chalt, err := pubkey.VerifyChallengeToken(chaltoken, "test", cookie)
	if err != nil {
		t.Fatalf("VerifyAuthenticationToken(%s) returns err=%+v", chaltoken, err)
	}

	claims := chalt.Claims.(jwt.MapClaims)
	if _, ok := claims["sub"]; !ok {
		t.Errorf("VerifyAuthenticationToken(%s) returns no sub", chaltoken)
	}

	if claims["sub"].(string) != "test" {
		t.Errorf("VerifyAuthenticationToken(%s) returns wrong sub: got=%s, want=%s", chaltoken, claims["sub"].(string), "test")
	}

	token, err := key.GenerateAuthenticationToken("test", "tester")
	if err != nil {
		t.Fatalf("GenerateAuthenticationToken(%s) returns err=%+v", "test", err)
	}
	t.Logf("GenerateAuthenticationToken(%s) returns token=%s", "test", token)

	parsed, err := pubkey.VerifyAuthenticationToken(token, "tester")
	if err != nil {
		t.Fatalf("VerifyAuthenticationToken(%s) returns err=%+v", token, err)
	}

	claims = parsed.Claims.(jwt.MapClaims)
	if _, ok := claims["sub"]; !ok {
		t.Errorf("VerifyAuthenticationToken(%s) returns no sub", token)
	}

	if claims["sub"].(string) != "test" {
		t.Errorf("VerifyAuthenticationToken(%s) returns wrong sub: got=%s, want=%s", token, claims["sub"].(string), "test")
	}
}
