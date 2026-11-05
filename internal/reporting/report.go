package reporting

import (
	"encoding/json"
	"fmt"
	"github.com/Diptesh2006/EnvGhost/pkg/model"
	"html/template"
	"os"
	"path/filepath"
)

func Save(r model.Run, dir string) error {
	if e := os.MkdirAll(dir, 0755); e != nil {
		return e
	}
	b, _ := json.MarshalIndent(r, "", "  ")
	if e := os.WriteFile(filepath.Join(dir, "latest.json"), b, 0644); e != nil {
		return e
	}
	const h = `<!doctype html><meta charset="utf-8"><title>EnvGhost report</title><h1>EnvGhost: {{.Config.Project.Name}}</h1><p>Baseline: {{.Baseline.Passed}} | Experiments: {{len .Experiments}}</p><table border="1"><tr><th>Dependency</th><th>Change</th><th>Classification</th><th>Logs</th></tr>{{range .Experiments}}<tr><td>{{.Name}}</td><td>{{.Change}}</td><td>{{.Classification}}</td><td><pre>{{.Execution.Logs}}</pre></td></tr>{{end}}</table>`
	t := template.Must(template.New("r").Parse(h))
	f, e := os.Create(filepath.Join(dir, "latest.html"))
	if e != nil {
		return e
	}
	defer f.Close()
	return t.Execute(f, r)
}
func Explain(r model.Run, name string) string {
	for _, x := range r.Experiments {
		if x.Name == name {
			return fmt.Sprintf("Dependency: %s\nChange: %s\nClassification: %s\nBaseline passed: %t\nExperiment passed: %t\nLogs:\n%s", x.Name, x.Change, x.Classification, x.Baseline.Passed, x.Execution.Passed, x.Execution.Logs)
		}
	}
	return "experiment not found"
}
