package authentication

import (
	"testing"

	"github.com/d1360-64rc14/simple-api/config"
)

var validSettings = &config.Auth{
	Base64TokenSeed: "IXRoZXF1aWNrZm94anVtcHNvdmVydGhlbGF6eWRvZyE",
	BCryptCost:      12,
}

func TestNewJWTEd25519Authenticator_WithInvalidSeedLength(t *testing.T) {
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

func TestNewJWTEd25519Authenticator_WithValidSeedLength(t *testing.T) {
	authenticator, err := NewJWTEd25519Authenticator(validSettings)
	if err != nil {
		t.Fatal(err)
	}
	if authenticator == nil {
		t.Error("authenticator was nil")
	}
}
