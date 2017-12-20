#!/usr/bin/sh

function get_tag () {
    tag=$(git describe --exact-match --tags)

    if ! [ $? -eq 0 ]; then
        tag='unknown'
    else
        tag=$(echo $tag | cut -d '-' -f 1,2)
    fi

    echo $tag
}

function get_branch () {
    branch=$(git rev-parse --abbrev-ref HEAD)

    if ! [ $? -eq 0 ]; then
        branch='unknown'
    fi

    echo $branch
}

function get_commit () {
    commit=$(git rev-parse HEAD)

    if ! [ $? -eq 0 ]; then
        commit='unknown'
    fi

    echo $commit
}


function main() {
    echo -e "\n========================================================"
    echo -e "start get version ..."

    TAG=$(get_tag)
    BRANCH=$(get_branch)
    COMMIT=$(get_commit)
    DATE=$(date '+%Y-%m-%d %H:%M:%S')
    
    Path="openauth/api/version"
    echo -e "collect project verion from git: tag:$TAG, data:$DATE, branch:$BRANCH, commit:$COMMIT"

    echo -e "start build ..."
    echo -e ""
    docker run --rm -e 'CGO_ENABLED=0' -e 'GOOS=linux' -e 'GOARCH=amd64' -v "$PWD":/go/src/openauth -w /go/src/openauth golang:1.9 go build -v -a -o openauth -ldflags "-X '$Path.GIT_TAG=${TAG}' -X '$Path.GIT_BRANCH=${BRANCH}' -X '$Path.GIT_COMMIT=${COMMIT}' -X '$Path.BUILD_TIME=${DATE}' -X '$Path.GO_VERSION=go1.9 linux/amd64'" cmd/openauthd/main.go
    echo -e ""

    echo -e "build completed!, the binaray file in this diretory"
    echo -e "========================================================\n"
}

main
