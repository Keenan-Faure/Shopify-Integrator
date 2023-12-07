#!/bin/bash

# source .env file
source ./../.env

# 1.) retrieve the command param (which is the directory to copy the file to)

if [[ -z "$1" ]]; then
    echo "command line parameter not found.";
    exit 404;
fi

# 2.) check if the file exists locally in the pwd directory
# if it does then we copy it to the Parent OS

if [[ -f "$pwd/$2" ]]; then
    # copy file to OS
    # docker cp command or whatever
else
    echo "file name $2 not found.";
    exit 404;
fi

exit;