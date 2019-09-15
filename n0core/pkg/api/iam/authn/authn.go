package authn

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

const TokenExpireMinutes = 30

var privateKey *rsa.PrivateKey

var publicKey interface{}

func init() {
	pubkey, ok := os.LookupEnv("N0CORE_AUTHN_JWT_PUBLIC_KEY_PATH")
	if ok {
		if err := LoadPublicKey(pubkey); err != nil {
			log.Printf("[CRITICAL] set valid the N0CORE_AUTHN_JWT_PUBLIC_KEY_PATH environment variable: %s", err.Error())
		}
	} else {
		log.Printf("[CRITICAL] set valid the N0CORE_AUTHN_JWT_PUBLIC_KEY_PATH environment variable")
	}

	key, ok := os.LookupEnv("N0CORE_AUTHN_JWT_PRIVATE_KEY_PATH")
	if ok {
		if err := LoadKey(key); err != nil {
			log.Printf("[CRITICAL] set valid the N0CORE_AUTHN_JWT_PRIVATE_KEY_PATH environment variable (Default: ./id_rsa): %s", err.Error())
		}
	} else {
		log.Printf("[CRITICAL] set valid the N0CORE_AUTHN_JWT_PRIVATE_KEY_PATH environment variable")
	}
}

func LoadKey(keyPath string) error {
	rawPubkey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return errors.Errorf("failed to read the JWT public key from %s", keyPath)
	}
	privateKeyBlock, _ := pem.Decode(rawPubkey)
	if privateKeyBlock == nil {
		return errors.Errorf("failed to decode public key")
	}
	if privateKeyBlock.Type != "RSA PRIVATE KEY" {
		return errors.Errorf("got wrong key type %s, want %s", privateKeyBlock.Type, "RSA PRIVATE KEY")
	}
	privateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return errors.Errorf("parse returns err=%s", err.Error())
	}

	return nil
}

func LoadPublicKey(pubkeyPath string) error {
	rawKey, err := ioutil.ReadFile(pubkeyPath)
	if err != nil {
		return errors.Errorf("failed to read the JWT public key from %s", pubkeyPath)
	}
	publicKeyBlock, _ := pem.Decode(rawKey)
	if publicKeyBlock == nil {
		return errors.Errorf("failed to decode public key")
	}
	if publicKeyBlock.Type != "PUBLIC KEY" {
		return errors.Errorf("got wrong key type %s, want %s", publicKeyBlock.Type, "RSA PRIVATE KEY")
	}
	publicKey, err = x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return errors.Errorf("parse returns err=%s", err.Error())
	}

	return nil
}

func GetConnectingAccountName(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", grpc.Errorf(codes.Unauthenticated, "")
	}

	auth := md["authorization"]
	token, err := VerifyToken(auth[1])
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	if _, ok := claims["subject"]; !ok {
		return "", grpc.Errorf(codes.Unauthenticated, "")
	}

	return claims["subject"].(string), nil
}

func GenerateToken(username string) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   username,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(TokenExpireMinutes * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			err := errors.New("Unexpected signing method")
			return nil, err
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "token is invalid")
	}
	if !parsedToken.Valid {
		return nil, errors.New("token is invalid")
	}

	return parsedToken, nil
}
