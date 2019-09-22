package jwtutil

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
	service := "test_project.example.com"
	var challengeToken string
	{
		challengeToken, _ = challengePrivateKey.GenerateChallengeToken(username, service, challenge)
	}

	// idp
	kg := NewKeyGenerator([]byte("secret"))
	authnPrivateKey, _ := kg.Generate(service)
	authnPublicKey, _ := authnPrivateKey.PublicKey()
	var authnToken string
	{
		if err := challengePublicKey.VerifyChallengeToken(challengeToken, username, service, challenge); err != nil {
			panic("failed verification")
		}

		authnToken, _ = authnPrivateKey.GenerateAuthenticationToken(username, service)
	}

	// service provider
	{
		serviceClient, err := authnPublicKey.VerifyAuthenticationToken(authnToken, service)
		if err != nil {
			panic("failed verification")
		}

		fmt.Printf("%s is using the service", serviceClient)
	}
}
