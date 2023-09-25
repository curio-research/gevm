#!/bin/bash


# first, check 
# if go exists and is the right version
if ! command -v go &> /dev/null; then
    echo "go isn't configured correctly or is not installed."
    exit 1
else
    v=`go version | { read _ _ v _; echo ${v#go}; }`

    if [[ "$v" != "1.20" ]]; then
        echo "the go version is incorrect."
        exit 1
    else
        # echo "go version: $v"
        cd examples/sum
        OUTPUT=$(go run main.go)
        echo "ran the main script, cleaning up files"
        rm -rf ./gevm-db/
        rm ./gevm
        echo $OUTPUT
    fi 
fi
