package config

import (
	"fmt"
	"github.com/Diptesh2006/EnvGhost/pkg/model"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

const File = ".envghost.yaml"

func Default() model.Config {
	return model.Config{Version: 1, Project: model.ProjectConfig{Name: "sample-api", Path: "."}, Runtime: model.RuntimeConfig{Image: "python:3.12-slim", Workdir: "/app", Setup: []string{"pip install -r requirements.txt"}, Command: "python main.py"}, Verification: model.VerificationConfig{StartupTimeoutText: "30s", CommandTimeoutText: "15s"}, Experiments: model.ExperimentConfig{MaxParallel: 1}, Output: model.OutputConfig{Directory: ".envghost/reports", Formats: []string{"json", "html"}}}
}
func Load(dir string) (model.Config, error) {
	b, e := os.ReadFile(filepath.Join(dir, File))
	if e != nil {
		return model.Config{}, fmt.Errorf("read %s: %w", File, e)
	}
	c := Default()
	if e = yaml.Unmarshal(b, &c); e != nil {
		return c, fmt.Errorf("invalid YAML: %w", e)
	}
	if e = Validate(&c); e != nil {
		return c, e
	}
	return c, nil
}
func Validate(c *model.Config) error {
	if c.Version != 1 {
		return fmt.Errorf("version must be 1")
	}
	if c.Runtime.Image == "" || c.Runtime.Command == "" {
		return fmt.Errorf("runtime.image and runtime.command are required")
	}
	if c.Project.Path == "" {
		c.Project.Path = "."
	}
	if c.Project.Name == "" {
		c.Project.Name = "project"
	}
	if c.Experiments.MaxParallel < 1 {
		c.Experiments.MaxParallel = 1
	}
	if c.Verification.StartupTimeoutText == "" {
		c.Verification.StartupTimeoutText = "30s"
	}
	if c.Verification.CommandTimeoutText == "" {
		c.Verification.CommandTimeoutText = "15s"
	}
	if _, e := time.ParseDuration(c.Verification.StartupTimeoutText); e != nil {
		return fmt.Errorf("verification.startup_timeout: %w", e)
	}
	return nil
}
func WriteDefault(dir string) error {
	b, _ := yaml.Marshal(Default())
	return os.WriteFile(filepath.Join(dir, File), b, 0644)
}
