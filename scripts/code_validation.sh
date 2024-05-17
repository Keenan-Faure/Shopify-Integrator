#!/bin/bash

# If you are unable to run this file then run
# chmod +x ./scripts/code_validation.sh

function start_app() {
    chmod +x ./scripts/install.github.sh
    ./scripts/install.github.sh
}

function go_test() {
    go test ./... -cover
}

function go_security() {
    go install github.com/securego/gosec/v2/cmd/gosec@latest
    gosec -exclude-generated ./...
}

function go_formatting() {
    test -z $(go fmt ./...)
}

function go_static_check() {
    go install honnef.co/go/tools/cmd/staticcheck@latest
    staticcheck ./...
}

function check_err() {
    if [ $? -eq 0 ]; then
        echo "$1: 👍 OK"
    else
        echo "$1: 👎 FAIL"
    fi
}

# start_app & BACK_PID=$!
# wait $BACK_PID

# go_test & BACK_PID=$!
# wait $BACK_PID

go_security &> /dev/null
check_err "Security"

go_formatting &> /dev/null
check_err "Formatting"

go_static_check &> /dev/null
check_err "Static Check"