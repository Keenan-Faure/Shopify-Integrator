#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/hosts.sh

function config_hosts() {
    source .env
    list="$APP_HOST $DOCS_HOST $API_SERVER_HOST"
    for item in $list
    do
        if ! append_host $item; then
            echo "error setting up /etc/hosts for $item"
            exit;
        fi
    done
}

function append_host() {
    if ! grep -q "$1" /etc/hosts; then
        echo "127.0.0.1   $1 >> /etc/hosts"
        echo "127.0.0.1   $1" >> /etc/hosts
    fi
}

config_hosts