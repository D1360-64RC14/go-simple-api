package authentication

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/d1360-64rc14/simple-api/config"
	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/golang-jwt/jwt/v5"
)

var _ interfaces.Authenticator = (*Ed25519JWTAuthenticator)(nil)

type Ed25519JWTAuthenticator struct {
	privKey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

func NewEd25519JWTAuthenticator(settings *config.Auth) (interfaces.Authenticator, error) {
	seed, err := base64.RawStdEncoding.DecodeString(settings.Base64TokenSeed)
	if err != nil {
		return nil, err
	}

	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey, ok := privateKey.Public().(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("could not parse crypto.PublicKey into an ed25519.PublicKey")
	}

	return &Ed25519JWTAuthenticator{
		privKey: privateKey,
		pubKey:  publicKey,
	}, nil
}

func (a Ed25519JWTAuthenticator) GenerateToken(id int, email string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodEdDSA,
		jwt.MapClaims{
			"id":    id,
			"email": email,
		},
	)

	return token.SignedString(a.privKey)
}

func (a Ed25519JWTAuthenticator) IsTokenValid(inputToken string) (bool, error) {
	token, err := jwt.Parse(inputToken, a.keyFunc)
	if err != nil {
		return false, err
	}

	return token.Valid, err
}

func (a Ed25519JWTAuthenticator) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
		return nil, fmt.Errorf("unexpected signing method %s", token.Method.Alg())
	}

	return a.privKey, nil
}
