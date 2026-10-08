package github

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	keyOnce sync.Once
	testKey *rsa.PrivateKey
)

// key gives one 2048-bit RSA key per test binary, made at run time: no key
// file enters the tree.
func key(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	keyOnce.Do(func() {
		k, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		testKey = k
	})
	return testKey
}

// checkJWT checks a JWT of the App (forge.md, The App identity) made at now
// and gives an error text, or "".
func checkJWT(jwt string, pub *rsa.PublicKey, appID int64, now time.Time) string {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return "the JWT has not three parts"
	}
	var header struct{ Alg, Typ string }
	var claims struct {
		Iat, Exp int64
		Iss      string
	}
	for i, v := range []any{&header, &claims} {
		raw, err := base64.RawURLEncoding.DecodeString(parts[i])
		if err != nil || json.Unmarshal(raw, v) != nil {
			return "a part of the JWT cannot be read"
		}
	}
	if header.Alg != "RS256" || header.Typ != "JWT" {
		return "the header is not RS256, JWT"
	}
	if claims.Iat != now.Unix()-60 {
		return "iat is not 60 seconds in the past"
	}
	if claims.Exp != claims.Iat+9*60 {
		return "exp is not iat + 9 minutes"
	}
	if claims.Iss != "42" || appID != 42 {
		return "iss is not the App ID"
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return "the signature cannot be read"
	}
	sum := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	if rsa.VerifyPKCS1v15(pub, crypto.SHA256, sum[:], sig) != nil {
		return "the signature does not verify with the public key"
	}
	return ""
}

func TestTheJWTHasItsHeaderClaimsAndSignature(t *testing.T) {
	k := key(t)
	now := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	jwt, err := makeJWT(k, 42, now)
	if err != nil {
		t.Fatal(err)
	}
	if msg := checkJWT(jwt, &k.PublicKey, 42, now); msg != "" {
		t.Error(msg)
	}
	other, _ := rsa.GenerateKey(rand.Reader, 2048)
	if checkJWT(jwt, &other.PublicKey, 42, now) == "" {
		t.Error("the signature verifies with another key")
	}
}

func TestFirstLineOfAMessage(t *testing.T) {
	for in, want := range map[string]string{
		`{"message":"Not Found","documentation_url":"x"}`: "Not Found",
		`{"message":"Bad\nmore"}`:                         "Bad",
		"<html>\n<body>":                                  "<html>",
		"":                                                "(no message)",
	} {
		if got := firstLine([]byte(in)); got != want {
			t.Errorf("firstLine(%q) = %q, want %q", in, got, want)
		}
	}
}
