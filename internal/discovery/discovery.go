package discovery

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

type Candidate struct{ Name, Source string }

var patterns = []struct {
	ext string
	re  *regexp.Regexp
}{{".py", regexp.MustCompile(`(?:os\.getenv\s*\(\s*["']|os\.environ(?:\.get\s*\(\s*["']|\s*\[\s*["']))([A-Z][A-Z0-9_]*)`)}, {".js", regexp.MustCompile(`process\.env\.([A-Z][A-Z0-9_]*)`)}}

func Scan(root string) []Candidate {
	m := map[string]string{}
	filepath.Walk(root, func(p string, i os.FileInfo, e error) error {
		if e != nil || i == nil || i.IsDir() {
			return nil
		}
		b, _ := os.ReadFile(p)
		for _, x := range patterns {
			if filepath.Ext(p) == x.ext {
				for _, q := range x.re.FindAllStringSubmatch(string(b), -1) {
					if len(q) > 1 {
						m[q[1]] = p
					}
				}
			}
		}
		return nil
	})
	out := []Candidate{}
	for n, s := range m {
		out = append(out, Candidate{n, s})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
