# EnvGhost

“It works on my machine. But why?” EnvGhost runs a configured application in Docker, establishes a passing baseline, then removes or changes one explicitly configured environmental assumption at a time. It reports observed regressions; static discovery is only a hint.

## Quick start

```sh
go build -o envghost ./cmd/envghost
cd examples/node-hidden-deps
../../envghost init # or use the checked-in configuration
../../envghost check . --max-parallel 2
../../envghost explain REQUIRED_TOKEN
```

Docker Desktop/Engine must be running. Applications run without host networking, privileges, Docker-socket access, or secret `.env` copying. The source is mounted read-only for the MVP; Docker is not a perfect hostile-code boundary.

## Commands

`envghost init` creates `.envghost.yaml`; `check [path]` runs the baseline and experiments; `report` prints JSON; `explain NAME` prints evidence; `doctor` checks Docker connectivity. Configure `runtime.environment` explicitly and use `experiments.environment` for unset/mutate experiments or `experiments.files` for file removal. The examples are reproducible with `go run ../../cmd/envghost check .` from their directories.

## Development

```sh
go test ./...
go test -race ./...
go vet ./...
go build ./cmd/envghost
```

The core supports Python and Node.js source discovery. HTTP verification, disposable service dependencies, richer file experiments, Docker API lifecycle control, and integration tests are planned next.

## License

MIT. See [LICENSE](LICENSE).
