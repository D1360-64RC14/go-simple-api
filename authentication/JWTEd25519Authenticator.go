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

var _ interfaces.Authenticator = (*JWTEd25519Authenticator)(nil)

type JWTEd25519Authenticator struct {
	privKey ed25519.PrivateKey
	pubKey  ed25519.PublicKey
}

func NewJWTEd25519Authenticator(settings *config.Auth) (interfaces.Authenticator, error) {
	seed, err := base64.RawStdEncoding.DecodeString(settings.Base64TokenSeed)
	if err != nil {
		return nil, err
	}

	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey, ok := privateKey.Public().(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("could not parse crypto.PublicKey into an ed25519.PublicKey")
	}

	return &JWTEd25519Authenticator{
		privKey: privateKey,
		pubKey:  publicKey,
	}, nil
}

func (a JWTEd25519Authenticator) GenerateToken(id int, email string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodEdDSA,
		jwt.MapClaims{
			"id":    id,
			"email": email,
		},
	)

	return token.SignedString(a.privKey)
}

func (a JWTEd25519Authenticator) IsTokenValid(inputToken string) (bool, error) {
	token, err := jwt.Parse(inputToken, a.keyFunc)
	if err != nil {
		return false, err
	}

	return token.Valid, err
}

func (a JWTEd25519Authenticator) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
		return nil, fmt.Errorf("unexpected signing method %s", token.Method.Alg())
	}

	return a.privKey, nil
}
