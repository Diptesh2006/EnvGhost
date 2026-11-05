package model

import "time"

type Config struct {
	Version      int                `yaml:"version"`
	Project      ProjectConfig      `yaml:"project"`
	Runtime      RuntimeConfig      `yaml:"runtime"`
	Verification VerificationConfig `yaml:"verification"`
	Experiments  ExperimentConfig   `yaml:"experiments"`
	Output       OutputConfig       `yaml:"output"`
}
type ProjectConfig struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}
type RuntimeConfig struct {
	Image       string            `yaml:"image"`
	Workdir     string            `yaml:"workdir"`
	Setup       []string          `yaml:"setup"`
	Command     string            `yaml:"command"`
	Environment map[string]string `yaml:"environment"`
}
type VerificationConfig struct {
	StartupTimeout     time.Duration `yaml:"-"`
	StartupTimeoutText string        `yaml:"startup_timeout"`
	CommandTimeout     time.Duration `yaml:"-"`
	CommandTimeoutText string        `yaml:"command_timeout"`
	HTTP               *HTTPConfig   `yaml:"http,omitempty"`
	SmokeCommand       string        `yaml:"smoke_command,omitempty"`
}
type HTTPConfig struct {
	Enabled        bool   `yaml:"enabled"`
	Port           int    `yaml:"port"`
	Path           string `yaml:"path"`
	ExpectedStatus int    `yaml:"expected_status"`
}
type ExperimentConfig struct {
	MaxParallel int                     `yaml:"max_parallel"`
	Environment []EnvironmentExperiment `yaml:"environment"`
	Files       []FileExperiment        `yaml:"files"`
}
type EnvironmentExperiment struct {
	Name   string `yaml:"name"`
	Action string `yaml:"action"`
	Value  string `yaml:"value,omitempty"`
}
type FileExperiment struct {
	Path   string `yaml:"path"`
	Action string `yaml:"action"`
}
type OutputConfig struct {
	Directory string   `yaml:"directory"`
	Formats   []string `yaml:"formats"`
}
type Run struct {
	ID          string             `json:"id"`
	StartedAt   time.Time          `json:"started_at"`
	Config      Config             `json:"config"`
	Baseline    Execution          `json:"baseline"`
	Experiments []ExperimentResult `json:"experiments"`
}
type Execution struct {
	Passed       bool          `json:"passed"`
	ExitCode     int           `json:"exit_code"`
	TimedOut     bool          `json:"timed_out"`
	Duration     time.Duration `json:"duration"`
	Logs         string        `json:"logs"`
	Verification string        `json:"verification"`
}
type ExperimentResult struct {
	Name           string    `json:"name"`
	Change         string    `json:"change"`
	Classification string    `json:"classification"`
	Execution      Execution `json:"execution"`
	Baseline       Execution `json:"baseline"`
}

const (
	Regression           = "REGRESSION"
	NoObservedRegression = "NO_OBSERVED_REGRESSION"
	Inconclusive         = "INCONCLUSIVE"
	ExperimentError      = "EXPERIMENT_ERROR"
)
