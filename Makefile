build:
	go build ./cmd/envghost
test:
	go test ./...
integration:
	go test -tags=integration ./tests/integration/...
lint:
	go vet ./...
demo:
	go run ./cmd/envghost check examples/node-hidden-deps
clean:
	go clean
