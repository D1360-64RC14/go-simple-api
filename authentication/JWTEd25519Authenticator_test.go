package authentication

import (
	"testing"

	"github.com/d1360-64rc14/simple-api/config"
	"github.com/d1360-64rc14/simple-api/interfaces"
)

var validSettings = &config.Auth{
	Base64TokenSeed: "IXRoZXF1aWNrZm94anVtcHNvdmVydGhlbGF6eWRvZyE",
	BCryptCost:      12,
}

func TestNewJWTEd25519Authenticator(t *testing.T) {
	t.Run("WithValidSeedLength", testNewJWTEd25519Authenticator_WithValidSeedLength)
	t.Run("WithInvalidSeedLength", testNewJWTEd25519Authenticator_WithInvalidSeedLength)
}

func testNewJWTEd25519Authenticator_WithValidSeedLength(t *testing.T) {
	authenticator, err := NewJWTEd25519Authenticator(validSettings)
	if err != nil {
		t.Fatal(err)
	}
	if authenticator == nil {
		t.Error("authenticator was nil")
	}
}

func testNewJWTEd25519Authenticator_WithInvalidSeedLength(t *testing.T) {
	invalidSettings := &config.Auth{
		Base64TokenSeed: "dGhlcXVpY2tmb3hqdW1wc292ZXJ0aGVsYXp5ZG9nIQ", // Token have 31 chars
		BCryptCost:      12,
	}

	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("should panic 'bad seed length'")
		}
		if rec != nil && rec.(string) != "ed25519: bad seed length: 31" {
			t.Fatalf("error should be 'ed25519: bad seed length: 31', got '%s'", rec.(string))
		}
	}()

	NewJWTEd25519Authenticator(invalidSettings)
}

func TestGenerateToken(t *testing.T) {
	authenticator, err := NewJWTEd25519Authenticator(validSettings)
	if err != nil {
		t.Fatal(err)
	}
	resultToken := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haWxAc2VydmVyLmNvbSIsImlkIjoxMjN9.L4CMw6zhZBrfPNs5QhPr3XebPqgKuf1ffi8QkYyK3WK9LoNN73p8bnt761PNykV4GOdJC3A3rqBnT33G1a6IBA"

	tokenStr, err := authenticator.GenerateToken(123, "mail@server.com")
	if err != nil {
		t.FailNow()
	}

	if tokenStr != resultToken {
		t.Errorf("token should be '%s', got '%s'", resultToken, tokenStr)
	}
}

func TestTokenIsValid(t *testing.T) {
	authenticator, err := NewJWTEd25519Authenticator(validSettings)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("WithValidToken", testTokenIsValid_WithValidToken(authenticator))
	t.Run("WithTamperedToken", testTokenIsValid_WithTamperedToken(authenticator))
}

func testTokenIsValid_WithValidToken(authenticator interfaces.Authenticator) func(*testing.T) {
	return func(t *testing.T) {
		validToken := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haWxAc2VydmVyLmNvbSIsImlkIjoxMjN9.L4CMw6zhZBrfPNs5QhPr3XebPqgKuf1ffi8QkYyK3WK9LoNN73p8bnt761PNykV4GOdJC3A3rqBnT33G1a6IBA"

		valid := authenticator.IsTokenValid(validToken)

		if !valid {
			t.Error("token should be valid, got invalid")
		}
	}
}

func testTokenIsValid_WithTamperedToken(authenticator interfaces.Authenticator) func(*testing.T) {
	return func(t *testing.T) {
		tamperedTokens := []string{
			"eyJhbGciOiJFRERTQSIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haWxAc2VydmVyLmNvbSIsImlkIjoxMjN9.L4CMw6zhZBrfPNs5QhPr3XebPqgKuf1ffi8QkYyK3WK9LoNN73p8bnt761PNykV4GOdJC3A3rqBnT33G1a6IBA",
			"eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im00MWxAc2VydmVyLmNvbSIsImlkIjoxMjN9.L4CMw6zhZBrfPNs5QhPr3XebPqgKuf1ffi8QkYyK3WK9LoNN73p8bnt761PNykV4GOdJC3A3rqBnT33G1a6IBA",
			"eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1haWxAc2VydmVyLmNvbSIsImlkIjoxMjN9.L4CMw6zhZBrfPNs5QhPr3XebPqgKuf1ffi8QkYyK3WK9LoNN73p8bnt761PNykV4GOdJC3A3rqBnT33G1a6IBa",
		}

		for _, tamperedToken := range tamperedTokens {
			valid := authenticator.IsTokenValid(tamperedToken)

			if valid {
				t.Error("token should be invalid, got valid")
			}
		}
	}
}
