#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./run.sh

OS="$(uname -s)"

# Builds the go code depending of OS
if [ $OS == "Darwin" ]; then
    echo "OSX detected"
    echo "GOOS=linux GOARCH=amd64 go build -o integrator"
    GOOS=linux GOARCH=amd64 go build -o integrator
else
    echo "Linux detected"
    echo "running go build -o integrator"
    go build -o integrator
fi

./integrator