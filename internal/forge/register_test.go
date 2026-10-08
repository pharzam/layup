package forge

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/fs"
	"math/big"
	"strings"
	"sync"
	"testing"
)

const registerHeader = "forge\tapp_id\tapp_slug\tkey_file\twatch_slug\tapi\tweb"

func register(rows ...string) []byte {
	return []byte(strings.Join(append([]string{registerHeader}, rows...), "\n") + "\n")
}

const goodRow = "github\t5118100\tlayup-agent\t/home/op/.config/layup-agent/app.pem\tlayup-watch\thttps://api.github.com\thttps://github.com"

func TestAGoodForgeRegisterIsRead(t *testing.T) {
	r, err := ReadForgeRegister(register(goodRow))
	if err != nil {
		t.Fatal(err)
	}
	want := Register{Forge: "github", AppID: 5118100, AppSlug: "layup-agent", KeyFile: "/home/op/.config/layup-agent/app.pem",
		WatchSlug: "layup-watch", API: "https://api.github.com", Web: "https://github.com"}
	if r != want {
		t.Fatalf("got %+v\nwant %+v", r, want)
	}
	if r, err := ReadForgeRegister(register(strings.Replace(goodRow, "\tlayup-watch\t", "\t—\t", 1))); err != nil || r.WatchSlug != "" {
		t.Errorf("an empty watch_slug: %+v, %v; want it read as empty", r, err)
	}
}

func TestForgeRegisterRefusesEachBrokenRule(t *testing.T) {
	for name, data := range map[string][]byte{
		"no row":                       register(),
		"a second row for github":      register(goodRow, goodRow),
		"an unknown column":            []byte(registerHeader + "\textra\n" + goodRow + "\tx\n"),
		"an app_id that is not an int": register(strings.Replace(goodRow, "5118100", "layup", 1)),
		"an empty app_slug":            register(strings.Replace(goodRow, "\tlayup-agent\t", "\t—\t", 1)),
		"an empty key_file":            register(strings.Replace(goodRow, "/home/op/.config/layup-agent/app.pem", "—", 1)),
		"a relative key_file":          register(strings.Replace(goodRow, "/home/op/.config/layup-agent/app.pem", "app.pem", 1)),
		"an empty api":                 register(strings.Replace(goodRow, "https://api.github.com", "—", 1)),
		"an empty web":                 register(strings.Replace(goodRow, "\thttps://github.com", "\t—", 1)),
	} {
		if r, err := ReadForgeRegister(data); err == nil {
			t.Errorf("%s: read %+v, want an error", name, r)
		} else if name != "no row" && !strings.Contains(err.Error(), "line ") {
			t.Errorf("%s: the error %q names no line", name, err)
		}
	}
}

var (
	keyOnce sync.Once
	testKey *rsa.PrivateKey
)

// key gives one 2048-bit RSA key per test binary, made at run time: no key
// file enters the tree (condition 7 of the plan review of #137).
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

func pemOf(typ string, der []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: typ, Bytes: der})
}

// pkcs1With gives the PKCS #1 DER of k with the version and the first prime
// that change gives it.
func pkcs1With(k *rsa.PrivateKey, version int, p *big.Int) []byte {
	type pkcs1 struct {
		Version             int
		N                   *big.Int
		E                   int
		D, P, Q, Dp, Dq, Qi *big.Int
	}
	der, err := asn1.Marshal(pkcs1{version, k.N, k.E, k.D, p, k.Primes[1], k.Precomputed.Dp, k.Precomputed.Dq, k.Precomputed.Qinv})
	if err != nil {
		panic(err)
	}
	return der
}

const runner = 501 // the ID of the user of the run, in the cases

func TestAGoodKeyIsChecked(t *testing.T) {
	k := key(t)
	pkcs8, err := x509.MarshalPKCS8PrivateKey(k)
	if err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string][]byte{
		"PKCS #1": pemOf("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(k)),
		"PKCS #8": pemOf("PRIVATE KEY", pkcs8),
	} {
		got, err := checkKey("/k.pem", 0o600, runner, runner, data)
		if err != nil {
			t.Errorf("%s: %v", name, err)
		} else if got == nil || !got.Equal(k) {
			t.Errorf("%s: the key read is not the key written", name)
		}
	}
}

func TestCheckKeyRefusesEachBrokenRule(t *testing.T) {
	k := key(t)
	good := pemOf("RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(k))
	der := x509.MarshalPKCS1PrivateKey(k)
	broken := append([]byte{}, der...)
	broken[0] ^= 0xff // not an ASN.1 SEQUENCE: the parse fails, not Validate
	otherAlgo, _ := asn1.Marshal(struct {
		Version int
		Algo    struct{ Algorithm asn1.ObjectIdentifier }
		Key     []byte
	}{0, struct{ Algorithm asn1.ObjectIdentifier }{asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}}, der})
	for _, c := range []struct {
		name         string
		mode         fs.FileMode
		owner        int
		data         []byte
		wantInReason string
	}{
		{"mode 0644", 0o644, runner, good, "0644"},
		{"mode 0400, not exactly 0600", 0o400, runner, good, "0400"},
		{"an owner other than the user of the run", 0o600, runner + 1, good, "owner"},
		{"no PEM block", 0o600, runner, []byte("not a key\n"), "PEM"},
		{"a PEM block of another type", 0o600, runner, pemOf("EC PRIVATE KEY", der), "EC PRIVATE KEY"},
		{"two PEM blocks", 0o600, runner, append(append([]byte{}, good...), good...), "one PEM block"},
		{"text before the PEM block", 0o600, runner, append([]byte("a note\n"), good...), "one PEM block"},
		{"a malformed DER", 0o600, runner, pemOf("RSA PRIVATE KEY", broken), "cannot be read"},
		{"a PKCS #1 key of version 1", 0o600, runner, pemOf("RSA PRIVATE KEY", pkcs1With(k, 1, k.Primes[0])), "version"},
		{"a key whose numbers fail Validate", 0o600, runner, pemOf("RSA PRIVATE KEY", pkcs1With(k, 0, new(big.Int).Add(k.Primes[0], big.NewInt(2)))), "RSA"},
		{"a PKCS #8 key of another algorithm", 0o600, runner, pemOf("PRIVATE KEY", otherAlgo), "algorithm"},
	} {
		_, err := checkKey("/k.pem", c.mode, c.owner, runner, c.data)
		switch {
		case err == nil:
			t.Errorf("%s: no error", c.name)
		case !strings.Contains(err.Error(), "/k.pem") || !strings.Contains(err.Error(), fmt.Sprintf("mode %04o", c.mode.Perm())):
			t.Errorf("%s: the error %q does not name the file and its mode (docs/spec/run.md, Input states)", c.name, err)
		case !strings.Contains(err.Error(), c.wantInReason):
			t.Errorf("%s: the error %q does not say %q", c.name, err, c.wantInReason)
		}
	}
}

func TestCheckKeyFileRefusesAMissingFile(t *testing.T) {
	path := t.TempDir() + "/missing.pem"
	if _, err := CheckKeyFile(path); err == nil || !strings.Contains(err.Error(), path) {
		t.Fatalf("CheckKeyFile(%s) = %v; want an error that names the file", path, err)
	}
}
