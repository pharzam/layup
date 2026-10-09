package session

import (
	"errors"
	"io/fs"
	"slices"
	"strings"
	"testing"
)

// files is a stand-in of the credential files: their content by path.
func files(m map[string]string) func(string) ([]byte, error) {
	return func(p string) ([]byte, error) {
		if c, ok := m[p]; ok {
			return []byte(c), nil
		}
		return nil, fs.ErrNotExist
	}
}

// hostile sets host variables that must never reach a session.
func hostile(t *testing.T) {
	t.Helper()
	t.Setenv("PATH", "/usr/local/bin:/usr/bin")
	for _, k := range []string{"GH_TOKEN", "GITHUB_TOKEN", "SSH_AUTH_SOCK", "HOME", "TMPDIR", "GIT_CONFIG_GLOBAL", "GIT_DIR", "ANTHROPIC_API_KEY"} {
		t.Setenv(k, "hostile")
	}
}

func TestTheEnvironmentIsTheNamedList(t *testing.T) {
	hostile(t)
	home, tmp := "/h/sessions/S-1a2b3c4d/home", "/h/sessions/S-1a2b3c4d/tmp"
	base := []string{"PATH=/usr/local/bin:/usr/bin", "LANG=C.UTF-8", "HOME=" + home, "TMPDIR=" + tmp, "GIT_CONFIG_NOSYSTEM=1"}
	read := files(map[string]string{"/k/one": "sk-one\n", "/k/two": "sk-two\n\n", "/k/none": "sk-none", "/k/cr": "sk-cr\r\n"})
	for name, c := range map[string]struct {
		h    Harness
		want []string
	}{
		"no credential, no variable":    {Harness{}, base},
		"var: with one final line feed": {Harness{Credential: "/k/one", CredentialTo: "var:ANTHROPIC_API_KEY"}, append(slices.Clone(base), "ANTHROPIC_API_KEY=sk-one")},
		"var: with two, one stays":      {Harness{Credential: "/k/two", CredentialTo: "var:KEY"}, append(slices.Clone(base), "KEY=sk-two\n")},
		"var: with none, unchanged":     {Harness{Credential: "/k/none", CredentialTo: "var:KEY"}, append(slices.Clone(base), "KEY=sk-none")},
		"var: a carriage return stays":  {Harness{Credential: "/k/cr", CredentialTo: "var:KEY"}, append(slices.Clone(base), "KEY=sk-cr\r")},
		"file: gives no variable":       {Harness{Credential: "/k/one", CredentialTo: "file:.config/devin/credentials.toml"}, base},
		"the fixed variables, in order": {Harness{Vars: []string{"ANTHROPIC_DEFAULT_HAIKU_MODEL=claude-opus-5-5", "A=b=c"}}, append(slices.Clone(base), "ANTHROPIC_DEFAULT_HAIKU_MODEL=claude-opus-5-5", "A=b=c")},
	} {
		got, err := environ(home, tmp, c.h, read)
		if err != nil || !slices.Equal(got, c.want) {
			t.Errorf("%s: %q, %v\nwant %q", name, got, err, c.want)
		}
		for _, kv := range got {
			if strings.Contains(kv, "hostile") {
				t.Errorf("%s: a host variable reached the session: %s", name, kv)
			}
		}
	}
	if _, err := environ(home, tmp, Harness{Credential: "/k/missing", CredentialTo: "var:KEY"}, read); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("a missing credential file: %v, want an error", err)
	}
}
