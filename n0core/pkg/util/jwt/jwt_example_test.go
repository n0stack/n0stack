package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"io"
)

func Example() {
	// idp
	challenge := make([]byte, 256)
	{
		io.ReadFull(rand.Reader, challenge[:])
	}

	// user
	challengeKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	challengePrivateKey, _ := NewPrivateKey(challengeKey)
	challengePublicKey, _ := challengePrivateKey.PublicKey()
	username := "test_user"
	var challengeToken string
	{
		challengeToken, _ = challengePrivateKey.GenerateChallengeToken(username, challenge)
	}

	// idp
	kg := NewKeyGenerator([]byte("secret"))
	project := "test_project"
	authnPrivateKey, authnPublicKey, _ := kg.Generate(project)
	var authnToken string
	{
		if err := challengePublicKey.VerifyChallengeToken(challengeToken, username, challenge); err != nil {
			panic("failed verification")
		}

		authnToken, _ = authnPrivateKey.GenerateAuthenticationToken(username, project)
	}

	// service provider
	{
		serviceClient, err := authnPublicKey.VerifyAuthenticationToken(authnToken, project)
		if err != nil {
			panic("failed verification")
		}

		fmt.Printf("%s is using the service", serviceClient)
	}
}
