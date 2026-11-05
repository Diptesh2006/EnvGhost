package reporting

import (
	"github.com/Diptesh2006/EnvGhost/pkg/model"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveWritesJSONAndEscapedHTML(t *testing.T) {
	d := t.TempDir()
	r := model.Run{Config: model.Config{Project: model.ProjectConfig{Name: "demo"}}, Experiments: []model.ExperimentResult{{Name: "X", Execution: model.Execution{Logs: "<script>"}}}}
	if err := Save(r, d); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"latest.json", "latest.html"} {
		if _, err := os.Stat(filepath.Join(d, name)); err != nil {
			t.Fatal(err)
		}
	}
	b, _ := os.ReadFile(filepath.Join(d, "latest.html"))
	if string(b) == "" || string(b) == "<script>" {
		t.Fatal("report was not rendered safely")
	}
}
