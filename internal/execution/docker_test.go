package execution

import (
	"context"
	"github.com/Diptesh2006/EnvGhost/pkg/model"
	"testing"
)

func TestRunMarksUnavailableDockerAsFailure(t *testing.T) {
	c := model.Config{Runtime: model.RuntimeConfig{Image: "does-not-matter", Workdir: "/app", Command: "true"}, Verification: model.VerificationConfig{CommandTimeoutText: "1s"}}
	r := Run(context.Background(), c, nil, "", "")
	if r.Passed {
		t.Fatal("expected unavailable Docker to fail")
	}
}
