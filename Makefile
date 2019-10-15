.PHONY: all build dist test install

build:
	@mkdir -p ./bin && rm -f ./bin/*
	go build -o ./bin/yaml ./cmd/yaml

dist:
	@mkdir -p ./bin && rm -f ./bin/*
	GOOS=darwin GOARCH=amd64 go build -o ./bin/yaml-darwin64 ./cmd/yaml
	GOOS=linux GOARCH=amd64 go build -o ./bin/yaml-linux64 ./cmd/yaml
	GOOS=linux GOARCH=386 go build -o ./bin/yaml-linux386 ./cmd/yaml
	GOOS=windows GOARCH=amd64 go build -o ./bin/yaml-windows64.exe ./cmd/yaml
	GOOS=windows GOARCH=386 go build -o ./bin/yaml-windows386.exe ./cmd/yaml

test:
	go test ./...

install:
	go install ./cmd/yaml
