// Package forge holds what LAYUP knows of the forge with no connection to it
// (docs/spec/forge.md; the table of M2a of docs/spec/packages.md): in row 22b
// of the plan (task T-1g1q, #137), the forge register of the host and the
// check of the App's key file. Row 24 adds the interface, its types, the
// capabilities and the permission check (O-171 of #137). It imports no
// network package (rule 5): the key is parsed with encoding/pem,
// encoding/asn1 and crypto/rsa, not crypto/x509, which depends on net.
package forge

import (
	"bytes"
	"crypto/rsa"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io/fs"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/pharzam/layup/internal/tsv"
)

// ForgeRegisterSchema is the form of host:registers/forge.tsv: the block
// forge-register of docs/spec/records.md.
var ForgeRegisterSchema = tsv.Schema{Name: "forge-register", Location: "host:registers/forge.tsv", Columns: []tsv.Column{
	{Name: "forge", Type: "enum(github)", Key: true},
	{Name: "app_id", Type: "int"},
	{Name: "app_slug", Type: "text"},
	{Name: "key_file", Type: "text"},
	{Name: "watch_slug", Type: "text"},
	{Name: "api", Type: "text"},
	{Name: "web", Type: "text"},
}}

// A Register is the one row of the forge register; WatchSlug is "" when the
// register holds the empty value.
type Register struct {
	Forge                                 string
	AppID                                 int64
	AppSlug, KeyFile, WatchSlug, API, Web string
}

// ReadForgeRegister reads the forge register by its schema and the rules of
// its block: one row (a register with no row is an error, docs/spec/run.md,
// Input states); each column but watch_slug holds a value; key_file is an
// absolute path. An error names its line.
func ReadForgeRegister(data []byte) (Register, error) {
	rows, err := tsv.Read(data, ForgeRegisterSchema)
	if err != nil {
		return Register{}, err
	}
	if len(rows) == 0 {
		return Register{}, &tsv.Error{Line: 1, Reason: "the forge register has no row"}
	}
	r := rows[0]
	for i, c := range ForgeRegisterSchema.Columns {
		if r[i] == "" && c.Name != "watch_slug" {
			return Register{}, &tsv.Error{Line: 2, Column: c.Name, Reason: "this column never holds the empty value"}
		}
	}
	if !filepath.IsAbs(r[3]) {
		return Register{}, &tsv.Error{Line: 2, Column: "key_file", Reason: "the key file is an absolute path"}
	}
	id, err := strconv.ParseInt(r[1], 10, 64)
	if err != nil {
		return Register{}, &tsv.Error{Line: 2, Column: "app_id", Reason: err.Error()}
	}
	return Register{r[0], id, r[2], r[3], r[4], r[5], r[6]}, nil
}

// CheckKeyFile checks the App's key file before any forge call
// (docs/spec/forge.md, The App identity) and gives its key: the file exists,
// its mode is exactly 0600, its owner is the user of the run, and it holds one
// PEM block of an RSA private key. Each error names the file. The owner is
// read from syscall.Stat_t, so the check runs on Unix hosts (macOS, the
// Linux of CI).
func CheckKeyFile(path string) (*rsa.PrivateKey, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("the key file %s: %w", path, err)
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("the key file %s: its owner cannot be read on this host", path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("the key file %s: %w", path, err)
	}
	return checkKey(path, info.Mode(), int(st.Uid), os.Getuid(), data)
}

// rsaOID is the algorithm of an RSA key in PKCS #8 (RFC 8017).
var rsaOID = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}

// checkKey is the check of CheckKeyFile on the mode, the owner, the ID of the
// user of the run and the bytes of the file, so each case runs on each host.
func checkKey(path string, mode fs.FileMode, owner, runner int, data []byte) (*rsa.PrivateKey, error) {
	bad := func(format string, a ...any) error {
		return fmt.Errorf("the key file %s: %s", path, fmt.Sprintf(format, a...))
	}
	if perm := mode.Perm(); perm != 0o600 || !mode.IsRegular() {
		return nil, bad("its mode is %04o, want 0600", perm)
	}
	if owner != runner {
		return nil, bad("its owner is the user %d, not the user of the run, %d", owner, runner)
	}
	block, rest := pem.Decode(data)
	if block == nil {
		return nil, bad("it holds no PEM block")
	}
	if len(bytes.TrimSpace(rest)) != 0 {
		return nil, bad("it holds more than one PEM block")
	}
	der := block.Bytes
	switch block.Type {
	case "RSA PRIVATE KEY":
	case "PRIVATE KEY":
		var p8 struct {
			Version int
			Algo    struct {
				Algorithm  asn1.ObjectIdentifier
				Parameters asn1.RawValue `asn1:"optional"`
			}
			Key []byte
		}
		if rest, err := asn1.Unmarshal(der, &p8); err != nil || len(rest) != 0 {
			return nil, bad("its PKCS #8 key cannot be read")
		}
		if !p8.Algo.Algorithm.Equal(rsaOID) {
			return nil, bad("its PKCS #8 key has the algorithm %v, not RSA", p8.Algo.Algorithm)
		}
		der = p8.Key
	default:
		return nil, bad("its PEM block is %s, not RSA PRIVATE KEY or PRIVATE KEY", block.Type)
	}
	var k struct {
		Version             int
		N                   *big.Int
		E                   int
		D, P, Q, Dp, Dq, Qi *big.Int
	}
	if rest, err := asn1.Unmarshal(der, &k); err != nil || len(rest) != 0 {
		return nil, bad("its RSA key cannot be read")
	}
	if k.Version != 0 {
		return nil, bad("its RSA key has the version %d, not 0 (two primes)", k.Version)
	}
	key := &rsa.PrivateKey{PublicKey: rsa.PublicKey{N: k.N, E: k.E}, D: k.D, Primes: []*big.Int{k.P, k.Q}}
	if err := key.Validate(); err != nil {
		return nil, bad("its RSA key is not valid: %v", err)
	}
	key.Precompute()
	return key, nil
}
