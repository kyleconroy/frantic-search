.PHONY: test fmt run release clean

gather: fmt gather.go
	go build

test: fmt
	go test -v

fmt:
	go fmt
