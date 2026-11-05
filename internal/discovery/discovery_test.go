package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanFindsPythonAndNodeCandidates(t *testing.T) {
	d := t.TempDir()
	if err := os.WriteFile(filepath.Join(d, "a.py"), []byte(`import os; os.getenv("DATABASE_URL")`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(d, "b.js"), []byte(`console.log(process.env.API_KEY)`), 0600); err != nil {
		t.Fatal(err)
	}
	got := Scan(d)
	if len(got) != 2 || got[0].Name != "API_KEY" || got[1].Name != "DATABASE_URL" {
		t.Fatalf("unexpected candidates: %#v", got)
	}
}
