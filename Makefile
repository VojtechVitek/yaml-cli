.PHONY: all build dist test install

build:
	@mkdir -p ./bin && rm -f ./bin/*
	go build -o ./bin/yaml ./cmd/yaml

dist:
	@mkdir -p ./bin && rm -f ./bin/*
	GOOS=darwin GOARCH=amd64 go build -gcflags=all="-l -B" -ldflags "-s -w" -o ./bin/yaml-darwin64 ./cmd/yaml
	GOOS=linux GOARCH=amd64 go build -gcflags=all="-l -B" -ldflags "-s -w" -o ./bin/yaml-linux-x64 ./cmd/yaml
	GOOS=windows GOARCH=amd64 go build -gcflags=all="-l -B" -ldflags "-s -w" -o ./bin/yaml-windows-x64.exe ./cmd/yaml
	GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ./bin/yaml-linux-x86 ./cmd/yaml
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -o ./bin/yaml-aarch64 ./cmd/yaml
	GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o ./bin/yaml-windows-x86.exe ./cmd/yaml
	GOOS=aix GOARCH=ppc64 go build -ldflags "-s -w" -o ./bin/yaml-aix-ppc64 ./cmd/yaml

compress:
	upx --best --ultra-brute bin/*

test:
	go test ./...

install:
	go install ./cmd/yaml
