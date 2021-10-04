.PHONY: serve
serve:
	go run ./cmd/server/main.go

.PHONY: test
test:
	go test -race ./...
