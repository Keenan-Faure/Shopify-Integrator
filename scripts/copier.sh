#!/bin/bash

# $1 -> the directory to which the file must be copied
# $2 -> the file name that should be copied

echo "preparing to copy..."

# source .env file
source ./../.env

# 1.) retrieve the command param (which is the directory to copy the file to)

if [[ -z "$1" ]]; then
    echo "command line parameter not found.";
    exit 404;
fi

# Checks if docker is install or running

if [ ! -x "$(command -v docker)" ]; then
    echo "Install docker"
    exit;
fi

# 2.) check if the file exists locally in the pwd directory
# if it does then we copy it to the Parent OS

var=$(pwd)

if [[ -f "$var/$2" ]]; then
    # copy file to OS
    if [ -d "$1" ]; then
        # checks if it's a docker container or local system
        if docker cp $SERVER_CONTAINER_NAME:$var/$2 $1 ; then
            exit;
        else
            # workaround is to copy it to the parent OS
            echo "running work-around"
            if cp $var/$2 $1 ; then
                exit;
            else
                # if none is able to copy then it errors
                echo "unable to copy :("
                exit 404;
            fi
        fi
    else 
        echo "invalid destination: $1";
        exit 404;
    fi
else
    echo "file name $var/$2 not found.";
    exit 404;
fi