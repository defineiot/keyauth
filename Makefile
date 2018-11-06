BINARY_NAME=keyauth
MAIN_FILE_PAHT=cmd/keyauth/main.go

all: test build

run:
		go build -o ${BINARY_NAME} ${MAIN_FILE_PAHT}
		./${BINARY_NAME} service bootstrap -f .keyauth/keyauth.conf

clean:
		go clean
		rm -f ./${BINARY_NAME}

test:
		go test -v ./...

build: local_build

linux_build:
		bash ./build/build.sh linux ${BINARY_NAME} ${MAIN_FILE_PAHT}

local_build:
		bash ./build/build.sh local ${BINARY_NAME} ${MAIN_FILE_PAHT}

docker_build:
		bash ./build/build.sh docker ${BINARY_NAME} ${MAIN_FILE_PAHT}


