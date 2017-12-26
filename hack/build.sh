#!/usr/bin/env bash

function _info(){
    local msg=$1
    local now=`date '+%Y-%m-%d %H:%M:%S'`
    echo -e "\033[44;37m [INFO] ${now} ${msg} \033[0m"
}
function get_tag () {
    local tag=$(git describe --exact-match --tags)

    if ! [ $? -eq 0 ]; then
        local tag='unknown'
    else
        local tag=$(echo ${tag} | cut -d '-' -f 1,2)
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


function main() {
    _info "start get version ..."

    TAG=$(get_tag)
    BRANCH=$(get_branch)
    COMMIT=$(get_commit)
    DATE=$(date '+%Y-%m-%d %H:%M:%S')

    Path="openauth/api/version"
    _info "collect project verion from git: tag:$TAG, data:$DATE, branch:$BRANCH, commit:$COMMIT"

    _info "start build ..."
    echo -e ""
    docker run --rm -e 'CGO_ENABLED=0' -e 'GOOS=linux' -e 'GOARCH=amd64' \
        -v "$PWD":/go/src/openauth \
        -w /go/src/openauth golang:1.9 \
        go build -v -a -o openauth -ldflags "-X '${Path}.GIT_TAG=${TAG}' -X '${Path}.GIT_BRANCH=${BRANCH}' -X '${Path}.GIT_COMMIT=${COMMIT}' -X '${Path}.BUILD_TIME=${DATE}' -X '${Path}.GO_VERSION=`go version`'" cmd/openauthd/main.go
    echo -e ""

    _info "build completed,the binary file in this directory."
}

main
