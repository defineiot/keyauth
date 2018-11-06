#!/usr/bin/env bash

BINARY_NAME=$2

function _info(){
    local msg=$1
    local now=`date '+%Y-%m-%d %H:%M:%S'`
    echo -e "\033[44;37m [INFO] ${now} ${msg} \033[0m"
}
function get_tag () {
    local tag=$(git describe --tags)

    if ! [ $? -eq 0 ]; then
        local tag='unknown'
    else
        local tag=$(echo ${tag} | cut -d '-' -f 1)
    fi

    echo ${tag}
}

function get_branch () {
    local branch=$(git rev-parse --abbrev-ref HEAD)

    if ! [ $? -eq 0 ]; then
        local branch='unknown'
    fi

    echo ${branch}
}

function get_commit () {
    local commit=$(git rev-parse HEAD)

    if ! [ $? -eq 0 ]; then
        local commit='unknown'
    fi

    echo ${commit}
}

function build () {
  local platform=$1
  local bin_name=$2
  local main_file=$3

  local version=$(go version | grep -o  'go[0-9].[0-9].*')

  echo ${platform}
  if [ ${platform} == "local" ]; then
    _info "start local build ..."
    echo -e ""
    go build -v -a -o ${bin_name} -ldflags "-X '${Path}.GIT_TAG=${TAG}' -X '${Path}.GIT_BRANCH=${BRANCH}' -X '${Path}.GIT_COMMIT=${COMMIT}' -X '${Path}.BUILD_TIME=${DATE}' -X '${Path}.GO_VERSION=${version}'" ${main_file}
    echo -e ""
  elif [ ${platform} == "linux" ]; then
     _info "start linux build ..."
    echo -e ""
    GOOS=linux GOARCH=amd64 \
        go build -v -a -o ${bin_name} -ldflags "-X '${Path}.GIT_TAG=${TAG}' -X '${Path}.GIT_BRANCH=${BRANCH}' -X '${Path}.GIT_COMMIT=${COMMIT}' -X '${Path}.BUILD_TIME=${DATE}' -X '${Path}.GO_VERSION=${version}'" ${main_file}
    echo -e ""
  elif [ ${platform} == "docker" ]; then
    _info "start docker build ..."
    echo -e ""
        docker run --rm -e 'CGO_ENABLED=0' -e 'GOOS=linux' -e 'GOARCH=amd64' \
        -v "$PWD":/go/src/iot-auth \
        -w /go/src/iot-auth golang:1.10.1 \
        go build -v -a -o ${bin_name} -ldflags "-X '${Path}.GIT_TAG=${TAG}' -X '${Path}.GIT_BRANCH=${BRANCH}' -X '${Path}.GIT_COMMIT=${COMMIT}' -X '${Path}.BUILD_TIME=${DATE}' -X '${Path}.GO_VERSION=${version}'" ${main_file}
    echo -e ""
  else
    echo "Please make sure the positon variable is local, docker or linux."
  fi
}

function main() {
    _info "start get version ..."

    TAG=$(get_tag)
    BRANCH=$(get_branch)
    COMMIT=$(get_commit)
    DATE=$(date '+%Y-%m-%d %H:%M:%S')

    Path="iot-auth/version"
    _info "collect project verion from git: tag:$TAG, data:$DATE, branch:$BRANCH, commit:$COMMIT"

    build $1 $2 $3

    _info "build completed,the binary file in this directory."
}

main $1 $2 $3
