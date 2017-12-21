#!/usr/bin/env bash

set -e

TOP_DIR=$(dirname `pwd`)

function find_unittest_package(){
    OLD_DIR=$(pwd)
    cd ${TOP_DIR}
    pkgs=$(find . -name *test.go | xargs dirname | uniq | xargs go list -e -f '{{.ImportPath}}')
    cd ${OLD_DIR}
    echo ${pkgs}
}
function run_unit_test(){
    pkgs=$(find_unittest_package)
    go test -cover ${pkgs}
}

run_unit_test