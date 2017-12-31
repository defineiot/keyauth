BINARY_NAME=openauthd

all: test build

run:
		go build -o ${BINARY_NAME} cmd/openauthd/main.go
		./${BINARY_NAME} service start

clean:
		go clean
		rm -f ./${BINARY_NAME}

test:
		go test -v ./...

build: build_in_local

build_in_docker:
		bash ./hack/build.sh

build_in_local:
		go build -o ${BINARY_NAME} cmd/openauthd/main.go
