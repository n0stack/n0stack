package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/hkdf"

	"golang.org/x/crypto/ssh"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

const AuthenticationTokenExpireMinutes = 30
const ChallengeTokenExpireMinutes = 3
const ChallengeTokenIssuer = "n0stack Challenge Response Authentication"

type PrivateKey struct {
	privateKey interface{}
	method     jwt.SigningMethod
}

func NewPrivateKey(key crypto.PrivateKey) (*PrivateKey, error) {
	privateKey := &PrivateKey{
		privateKey: key,
	}

	switch key.(type) {
	case *rsa.PrivateKey:
		privateKey.method = jwt.SigningMethodRS256
	case *ecdsa.PrivateKey:
		privateKey.method = jwt.SigningMethodES256
	default:
		return nil, errors.Errorf("unexpected key type")
	}

	return privateKey, nil
}

func ParsePrivateKey(in []byte) (*PrivateKey, error) {
	key, err := ssh.ParseRawPrivateKey(in)
	if err != nil {
		return nil, errors.Errorf("ParseRawPrivateKey() returns err=%s", err.Error())
	}

	return NewPrivateKey(key)
}

func ParsePrivateKeyFromFile(filename string) (*PrivateKey, error) {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", filename)
	}

	return ParsePrivateKey(in)
}

func (p PrivateKey) PublicKey() (*PublicKey, error) {
	type genpubkey interface {
		Public() crypto.PublicKey
	}

	genpub, ok := p.privateKey.(genpubkey)
	if !ok {
		return nil, errors.Errorf("unexpected private key")
	}

	pubkey := genpub.Public()
	return NewPublicKey(pubkey)
}

func (p PrivateKey) GenerateChallengeToken(username string, cookie []byte) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   username,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(AuthenticationTokenExpireMinutes * time.Minute).Unix(),
		Issuer:    ChallengeTokenIssuer,
		Id:        hex.EncodeToString(cookie[:]),
	}
	token := jwt.NewWithClaims(p.method, claims)
	t, err := token.SignedString(p.privateKey)

	return t, err
}

func (p PrivateKey) GenerateAuthenticationToken(username string, issuer string) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   username,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(AuthenticationTokenExpireMinutes * time.Minute).Unix(),
		Issuer:    issuer,
	}
	token := jwt.NewWithClaims(p.method, claims)
	t, err := token.SignedString(p.privateKey)

	return t, err
}

type PublicKey struct {
	publicKey interface{}
	method    jwt.SigningMethod
}

func NewPublicKey(key crypto.PublicKey) (*PublicKey, error) {
	pubkey := &PublicKey{
		publicKey: key,
	}

	switch pubkey.publicKey.(type) {
	case *rsa.PublicKey:
		pubkey.method = jwt.SigningMethodRS256
	case *ecdsa.PublicKey:
		pubkey.method = jwt.SigningMethodES256
	default:
		return nil, errors.Errorf("unexpected key type")
	}

	return pubkey, nil
}

func ParsePublicKey(in []byte) (*PublicKey, error) {
	var key interface{}
	var err error

	publicKeyBlock, _ := pem.Decode(in)
	if publicKeyBlock == nil {
		key, _, _, _, err = ssh.ParseAuthorizedKey(in)
		if err != nil {
			return nil, errors.Errorf("ParsePublicKey() returns err=%s", err.Error())
		}
	} else {
		if publicKeyBlock.Type != "PUBLIC KEY" {
			return nil, errors.Errorf("got wrong key type %s, want %s", publicKeyBlock.Type, "RSA PRIVATE KEY")
		}

		key, err = x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
		if err != nil {
			return nil, errors.Errorf("ParsePKIXPublicKey() returns err=%s", err.Error())
		}
	}

	type sshPublicKey interface {
		CryptoPublicKey() crypto.PublicKey
	}
	if k, ok := key.(sshPublicKey); ok {
		key = k.CryptoPublicKey()
	}

	return NewPublicKey(key)
}

func ParsePublicKeyFromFile(filename string) (*PublicKey, error) {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", filename)
	}

	return ParsePublicKey(in)
}

func (p PublicKey) VerifyChallengeToken(token string, username string, cookie []byte) error {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != p.method.Alg() {
			return nil, errors.Errorf("unexpected JWT algorithm: got=%s, want=%s", token.Method.Alg(), p.method.Alg())
		}

		claims := token.Claims.(jwt.MapClaims)
		if _, ok := claims["jti"]; !ok {
			return nil, errors.Errorf("token have no jti")
		}

		if jti := claims["jti"].(string); jti != hex.EncodeToString(cookie[:]) {
			return nil, errors.Errorf("failed cookie verification: got=%+v, want=%+v", jti, hex.EncodeToString(cookie[:]))
		}
		if iss := claims["iss"].(string); iss != ChallengeTokenIssuer {
			return nil, errors.Errorf("failed issuer verification: got=%s, want=%s", iss, ChallengeTokenIssuer)
		}
		if sub := claims["sub"].(string); sub != username {
			return nil, errors.Errorf("failed issuer verification: got=%s, want=%s", sub, username)
		}

		return p.publicKey, nil
	})

	if err != nil {
		return errors.Wrap(err, "token is invalid")
	}
	if !parsedToken.Valid {
		return errors.New("token is invalid")
	}

	return nil
}

func (p PublicKey) VerifyAuthenticationToken(token string, issuer string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != p.method.Alg() {
			return nil, errors.Errorf("unexpected JWT algorithm: got=%s, want=%s", token.Method.Alg(), p.method.Alg())
		}

		claims := token.Claims.(jwt.MapClaims)
		if iss := claims["iss"].(string); iss != issuer {
			return nil, errors.Errorf("failed issuer verification: got=%s, want=%s", iss, issuer)
		}

		return p.publicKey, nil
	})

	if err != nil {
		return "", errors.Wrap(err, "token is invalid")
	}
	if !parsedToken.Valid {
		return "", errors.New("token is invalid")
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	return claims["sub"].(string), nil
}

type KeyGenerator struct {
	secret []byte
}

func NewKeyGenerator(secret []byte) *KeyGenerator {
	return &KeyGenerator{
		secret: secret,
	}
}

func (k KeyGenerator) Generate(issuer string) (*PrivateKey, *PublicKey, error) {
	// Underlying hash function for HMAC.
	hash := sha256.New

	// Cryptographically secure master secret.
	// secret := []byte{0x00, 0x01, 0x02, 0x03} // i.e. NOT this.

	// Non-secret salt, optional (can be nil).
	// Recommended: hash-length random value.
	// salt := make([]byte, hash().Size())
	// if _, err := rand.Read(salt); err != nil {
	// 	panic(err)
	// }

	// Non-secret context info, optional (can be nil).
	info := []byte(issuer)

	// Generate three 128-bit derived keys.
	// nonceを保存するのは面倒
	kdf := hkdf.New(hash, k.secret, nil, info)

	key, err := ecdsa.GenerateKey(elliptic.P256(), kdf)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "ecdsa.GenerateKey() returns")
	}

	privkey, err := NewPrivateKey(key)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "NewPrivateKey() returns")
	}
	pubkey, err := privkey.PublicKey()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "PublicKey() returns")
	}

	return privkey, pubkey, nil
}
