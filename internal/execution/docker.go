package execution

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Diptesh2006/EnvGhost/pkg/model"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func DockerAvailable() error {
	c := exec.Command("docker", "info")
	if e := c.Run(); e != nil {
		return fmt.Errorf("EnvGhost could not connect to Docker. Check that Docker Desktop/Engine is running and try docker info: %w", e)
	}
	return nil
}
func Run(ctx context.Context, cfg model.Config, env map[string]string, remove string, removeFile string) model.Execution {
	start := time.Now()
	timeout, err := time.ParseDuration(cfg.Verification.CommandTimeoutText)
	if err != nil || timeout <= 0 {
		timeout = 15 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	args := []string{"run", "--rm", "--network", "none", "--memory", "512m", "--cpus", "1", "--pids-limit", "128", "-w", cfg.Runtime.Workdir}
	project, _ := filepath.Abs(cfg.Project.Path)
	args = append(args, "-v", project+":"+cfg.Runtime.Workdir+":ro")
	for k, v := range cfg.Runtime.Environment {
		if k != remove {
			args = append(args, "-e", k+"="+v)
		}
	}
	for k, v := range env {
		args = append(args, "-e", k+"="+v)
	}
	command := strings.Join(append(cfg.Runtime.Setup, cfg.Runtime.Command), " && ")
	if removeFile != "" {
		command = "rm -f -- " + shellQuote(removeFile) + " && " + command
	}
	args = append(args, cfg.Runtime.Image, "sh", "-c", command)
	cmd := exec.CommandContext(ctx, "docker", args...)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	e := cmd.Run()
	ex := model.Execution{Passed: e == nil, Duration: time.Since(start), Logs: b.String()}
	if ctx.Err() == context.DeadlineExceeded {
		ex.TimedOut = true
	}
	if e != nil {
		ex.ExitCode = 1
	}
	return ex
}

func shellQuote(value string) string { return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'" }
