package psb

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// golden runs Check on testdata/<name>.md and compares the TSV byte for byte
// with testdata/<name>.tsv. One mechanism, one test, every rule (issue #37).
func golden(t *testing.T, input, want string) {
	t.Helper()
	src, err := os.ReadFile(input)
	if err != nil {
		t.Fatal(err)
	}
	var got bytes.Buffer
	WriteTSV(&got, Check(string(src)))
	exp, err := os.ReadFile(want)
	if err != nil {
		t.Fatal(err)
	}
	if got.String() != string(exp) {
		t.Fatalf("%s: output differs from %s\n--- got ---\n%s--- want ---\n%s", input, want, got.String(), exp)
	}
}

func TestGolden(t *testing.T) {
	for _, name := range []string{"triggers", "clean", "edge"} {
		t.Run(name, func(t *testing.T) {
			golden(t, filepath.Join("testdata", name+".md"), filepath.Join("testdata", name+".tsv"))
		})
	}
}
