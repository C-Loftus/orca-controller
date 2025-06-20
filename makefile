
.PHONY: all gen test check

gen:
	go generate ./...

test: gen
	go test ./... -count 1

lint:
	golangci-lint run
