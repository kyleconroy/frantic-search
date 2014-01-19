.PHONY: test fmt run release clean

gather: fmt gather.go
	go build -o gather

test: fmt
	go test

fmt:
	go fmt

release: gather-osx.zip gather-linux.tar.gz

gather-osx.zip:
	rm -f gather
	GOOS=darwin GOARCH=386 go build
	zip -q gather-osx.zip gather

gather-linux.tar.gz:
	rm -f gather
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build
	tar czf gather-linux.tar.gz gather

gather-windows.zip:
	rm -f gather
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build
	zip -q gather-windows.zip gather


clean:
	rm -f gather
	rm -f gather-osx.zip
	rm -f gather-linux.tar.gz
