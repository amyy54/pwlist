all: test build

GIT_DESCRIBE=$(shell git describe)
GIT_DESCRIBE_LONG=$(shell git describe --long)

clean:
	rm -rf ./bin

test:
	go test ./...

build:
	go build -ldflags "-X 'main.Version=$(GIT_DESCRIBE)' -X 'main.VersionLong=$(GIT_DESCRIBE_LONG)'" -o bin/pwlist ./cmd/pwlist
