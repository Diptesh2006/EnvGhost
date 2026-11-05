package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Diptesh2006/EnvGhost/internal/config"
	"github.com/Diptesh2006/EnvGhost/internal/execution"
	"github.com/Diptesh2006/EnvGhost/internal/reporting"
	"github.com/Diptesh2006/EnvGhost/pkg/model"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var root = &cobra.Command{Use: "envghost"}

func init() { root.AddCommand(initCmd(), doctorCmd(), checkCmd(), reportCmd(), explainCmd()) }
func Execute() {
	if e := root.Execute(); e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
}
func initCmd() *cobra.Command {
	return &cobra.Command{Use: "init", RunE: func(cmd *cobra.Command, _ []string) error {
		if _, e := os.Stat(config.File); e == nil {
			return fmt.Errorf("%s already exists", config.File)
		}
		return config.WriteDefault(".")
	}}
}
func doctorCmd() *cobra.Command {
	return &cobra.Command{Use: "doctor", RunE: func(*cobra.Command, []string) error { return execution.DockerAvailable() }}
}
func checkCmd() *cobra.Command {
	var parallel int
	cmd := &cobra.Command{Use: "check [project]", Args: cobra.MaximumNArgs(1), RunE: func(_ *cobra.Command, a []string) error {
		dir := "."
		if len(a) > 0 {
			dir = a[0]
		}
		c, e := config.Load(dir)
		if e != nil {
			return e
		}
		if parallel > 0 {
			c.Experiments.MaxParallel = parallel
		}
		if e = execution.DockerAvailable(); e != nil {
			return e
		}
		r := model.Run{ID: "latest", Config: c}
		ctx := context.Background()
		r.Baseline = execution.Run(ctx, c, nil, "", "")
		if !r.Baseline.Passed {
			return fmt.Errorf("baseline failed; no dependency conclusions can be established\n%s", r.Baseline.Logs)
		}
		jobs := make(chan model.EnvironmentExperiment)
		out := make(chan model.ExperimentResult)
		var wg sync.WaitGroup
		for i := 0; i < c.Experiments.MaxParallel; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for x := range jobs {
					removeFile := ""
					if strings.HasPrefix(x.Action, "file ") {
						removeFile = x.Name
					}
					e := execution.Run(ctx, c, experimentEnv(x), unsetName(x), removeFile)
					class := model.NoObservedRegression
					if !e.Passed {
						class = model.Regression
					}
					out <- model.ExperimentResult{Name: x.Name, Change: x.Action, Classification: class, Execution: e, Baseline: r.Baseline}
				}
			}()
		}
		go func() {
			for _, x := range c.Experiments.Environment {
				jobs <- x
			}
			close(jobs)
			wg.Wait()
			close(out)
		}()
		for _, x := range c.Experiments.Files {
			jobs <- model.EnvironmentExperiment{Name: x.Path, Action: "file " + x.Action}
		}
		for x := range out {
			r.Experiments = append(r.Experiments, x)
		}
		od := filepath.Join(dir, c.Output.Directory)
		if e = reporting.Save(r, od); e != nil {
			return e
		}
		fmt.Printf("Baseline: PASS\nExperiments: %d\nReport: %s\n", len(r.Experiments), filepath.Join(od, "latest.html"))
		return nil
	}}
	cmd.Flags().IntVar(&parallel, "max-parallel", 0, "maximum concurrent experiments")
	return cmd
}

func experimentEnv(x model.EnvironmentExperiment) map[string]string {
	if x.Action == "mutate" {
		return map[string]string{x.Name: x.Value}
	}
	return map[string]string{}
}

func unsetName(x model.EnvironmentExperiment) string {
	if x.Action == "unset" {
		return x.Name
	}
	return ""
}
func reportCmd() *cobra.Command {
	return &cobra.Command{Use: "report", RunE: func(*cobra.Command, []string) error {
		b, e := os.ReadFile(".envghost/reports/latest.json")
		if e != nil {
			return e
		}
		fmt.Println(string(b))
		return nil
	}}
}
func explainCmd() *cobra.Command {
	return &cobra.Command{Use: "explain NAME", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, a []string) error {
		b, e := os.ReadFile(".envghost/reports/latest.json")
		if e != nil {
			return e
		}
		var r model.Run
		if e = json.Unmarshal(b, &r); e != nil {
			return e
		}
		fmt.Println(reporting.Explain(r, a[0]))
		return nil
	}}
}
