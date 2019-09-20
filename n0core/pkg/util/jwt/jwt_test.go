package jwt

import (
	"crypto/rand"
	"io"
	"testing"
)

// Key generation:
//   openssl genrsa 4096 > private.key
//   openssl rsa -pubout < private.key > public.key
func TestJWTWithOpenSSLRSA(t *testing.T) {
	parseFile(t, "private.key", "public.key")
}

// ssh-keygen -t ecdsa
func TestJWTWithSSHECDSA(t *testing.T) {
	parseFile(t, "id_ecdsa", "id_ecdsa.pub")
}

// ssh-keygen -t rsa
func TestJWTWithSSHRSA(t *testing.T) {
	parseFile(t, "id_rsa", "id_rsa.pub")
}

func TestGenerator(t *testing.T) {
	kg := NewKeyGenerator([]byte("foo"))
	key, pubkey, err := kg.Generate("bar")
	if err != nil {
		t.Fatalf("Generate(%s) returns err=%+v", "bar", err)
	}

	testAuthentication(t, key, pubkey)
}

func parseFile(t *testing.T, keyFile, pubkeyFile string) {
	key, err := ParsePrivateKeyFromFile(keyFile)
	if err != nil {
		t.Fatalf("ParsePrivateKeyFromFile() returns err=%+v", err)
	}
	pubkey, err := ParsePublicKeyFromFile(pubkeyFile)
	if err != nil {
		t.Fatalf("ParsePublicKeyFromFile() returns err=%+v", err)
	}

	testAuthentication(t, key, pubkey)
}

func testAuthentication(t *testing.T, key *PrivateKey, pubkey *PublicKey) {
	cookie := make([]byte, 256)
	io.ReadFull(rand.Reader, cookie[:])
	chaltoken, err := key.GenerateChallengeToken("test", cookie)
	if err != nil {
		t.Fatalf("GenerateAuthenticationToken(%s) returns err=%+v", "test", err)
	}
	t.Logf("GenerateAuthenticationToken(%s) returns token=%s", "test", chaltoken)

	if err := pubkey.VerifyChallengeToken(chaltoken, "test", cookie); err != nil {
		t.Fatalf("VerifyAuthenticationToken(%s) returns err=%+v", chaltoken, err)
	}

	token, err := key.GenerateAuthenticationToken("test", "tester")
	if err != nil {
		t.Fatalf("GenerateAuthenticationToken(%s) returns err=%+v", "test", err)
	}
	t.Logf("GenerateAuthenticationToken(%s) returns token=%s", "test", token)

	username, err := pubkey.VerifyAuthenticationToken(token, "tester")
	if err != nil {
		t.Fatalf("VerifyAuthenticationToken(%s) returns err=%+v", token, err)
	}

	if username != "test" {
		t.Errorf("VerifyAuthenticationToken(%s) returns wrong sub: got=%s, want=%s", token, username, "test")
	}
}

func Benchmark(b *testing.B) {
	keyFile := "id_ecdsa"
	pubkeyFile := "id_ecdsa.pub"
	key, err := ParsePrivateKeyFromFile(keyFile)
	if err != nil {
		b.Fatalf("ParsePrivateKeyFromFile() returns err=%+v", err)
	}
	pubkey, err := ParsePublicKeyFromFile(pubkeyFile)
	if err != nil {
		b.Fatalf("ParsePublicKeyFromFile() returns err=%+v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cookie := make([]byte, 256)
		io.ReadFull(rand.Reader, cookie[:])

		chaltoken, err := key.GenerateChallengeToken("test", cookie)
		if err != nil {
			b.Fatalf("GenerateAuthenticationToken(%s) returns err=%+v", "test", err)
		}

		if err := pubkey.VerifyChallengeToken(chaltoken, "test", cookie); err != nil {
			b.Fatalf("VerifyAuthenticationToken(%s) returns err=%+v", chaltoken, err)
		}

		token, err := key.GenerateAuthenticationToken("test", "tester")
		if err != nil {
			b.Fatalf("GenerateAuthenticationToken(%s) returns err=%+v", "test", err)
		}

		_, err = pubkey.VerifyAuthenticationToken(token, "tester")
		if err != nil {
			b.Fatalf("VerifyAuthenticationToken(%s) returns err=%+v", token, err)
		}
	}
}
