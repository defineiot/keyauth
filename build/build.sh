#!/usr/bin/env bash

BINARY_NAME=$2

function _info(){
    local msg=$1
    local now=`date '+%Y-%m-%d %H:%M:%S'`
    echo -e "\033[1;46;30m[INFO]\033[0m ${now} ${msg}"
}

function _version(){
    local msg=$1
    local now=`date '+%Y-%m-%d %H:%M:%S'`
    echo -e "\033[1;46;30m[INFO]\033[0m ${now} ${msg}"
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


  if [ ${platform} == "local" ]; then
    _info "开始本地构建 ..."
    echo -e ""
    go build -a -o ${bin_name} -ldflags "-X '${Path}.GIT_TAG=${TAG}' -X '${Path}.GIT_BRANCH=${BRANCH}' -X '${Path}.GIT_COMMIT=${COMMIT}' -X '${Path}.BUILD_TIME=${DATE}' -X '${Path}.GO_VERSION=${version}'" ${main_file}
    echo -e ""
  elif [ ${platform} == "linux" ]; then
     _info "开始构建Linux平台版本 ..."
    echo -e ""
    GOOS=linux GOARCH=amd64 \
        go build -a -o ${bin_name} -ldflags "-X '${Path}.GIT_TAG=${TAG}' -X '${Path}.GIT_BRANCH=${BRANCH}' -X '${Path}.GIT_COMMIT=${COMMIT}' -X '${Path}.BUILD_TIME=${DATE}' -X '${Path}.GO_VERSION=${version}'" ${main_file}
    echo -e ""
  elif [ ${platform} == "docker" ]; then
    _info "开始基于Docker ..."
    echo -e ""
        docker run --rm -e 'CGO_ENABLED=0' -e 'GOOS=linux' -e 'GOARCH=amd64' \
        -v "$PWD":/go/src/github.com/defineiot/keyauth \
        -w /go/src/github.com/defineiot/keyauth golang:1.10.1 \
        go build -v -a -o ${bin_name} -ldflags "-X '${Path}.GIT_TAG=${TAG}' -X '${Path}.GIT_BRANCH=${BRANCH}' -X '${Path}.GIT_COMMIT=${COMMIT}' -X '${Path}.BUILD_TIME=${DATE}' -X '${Path}.GO_VERSION=${version}'" ${main_file}
    echo -e ""
  else
    echo "Please make sure the positon variable is local, docker or linux."
  fi
}

function main() {
    _info "开始构建 [keyauth] ..."

    TAG=$(get_tag)
    BRANCH=$(get_branch)
    COMMIT=$(get_commit)
    DATE=$(date '+%Y-%m-%d %H:%M:%S')

    Path="github.com/defineiot/keyauth/version"
    _version "构建版本的时间(Build Time): $DATE"
    _version "当前构建的版本(Git   Tag ): $TAG"
    _version "当前构建的分支(Git Branch): $BRANCH"
    _version "当前构建的提交(Git Commit): $COMMIT"

    build $1 $2 $3

    _info "程序构建完成: $2"
}

main $1 $2 $3
