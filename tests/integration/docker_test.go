//go:build integration

package integration

import (
	"context"
	"github.com/Diptesh2006/EnvGhost/internal/execution"
	"github.com/Diptesh2006/EnvGhost/pkg/model"
	"testing"
)

func TestDockerExecution(t *testing.T) {
	if err := execution.DockerAvailable(); err != nil {
		t.Skip(err)
	}
	c := model.Config{Project: model.ProjectConfig{Path: "../../examples/node-hidden-deps"}, Runtime: model.RuntimeConfig{Image: "node:22-alpine", Workdir: "/app", Command: "node main.js", Environment: map[string]string{"REQUIRED_TOKEN": "demo"}}, Verification: model.VerificationConfig{CommandTimeoutText: "10s"}}
	r := execution.Run(context.Background(), c, nil, "", "")
	if !r.Passed {
		t.Fatalf("baseline failed: %s", r.Logs)
	}
	r = execution.Run(context.Background(), c, nil, "REQUIRED_TOKEN", "")
	if r.Passed {
		t.Fatal("unset required variable did not regress")
	}
}
