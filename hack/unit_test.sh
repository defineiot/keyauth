#!/usr/bin/env bash

set -e

function find_unittest_package(){
    pkgs=$(find . -name *test.go | xargs dirname | uniq | xargs go list -e -f '{{.ImportPath}}')
    echo ${pkgs}
}
function run_unit_test(){
    pkgs=$(find_unittest_package)
    go test -cover ${pkgs}
}

run_unit_test