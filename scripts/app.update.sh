#!/bin/bash

# Please do not modify this file, modify the .env file within this directory
# If you are unable to run this file then run
# chmod +x ./scripts/app.update.sh

function err_msg() {
    echo "invalid argument"
    echo "either enter:"
    echo "./scripts/app.update.sh ./app/{{file_name_with_extension}}  -  relative"
    echo "./scripts/app.update.sh $pwd/{{file_name_with_extension}} -  absolute"
    exit 1;
}

function cp_file() {
    if [[ ! $# -eq 0 ]] ; then
        if [ ! -z "$1" ]; then
            if [ -f $1 ]; then
                source .env
                echo "$APP_CONTAINER_NAME"
                if docker cp $1 $APP_CONTAINER_NAME:/keenan/$1 ; then
                    docker exec $APP_CONTAINER_NAME /bin/sh -c "cd /keenan/app/ && npm install"
                    docker restart $APP_CONTAINER_NAME
                else
                    exit;
                fi
            else
                echo "File at location '$1' does not exist."
                exit;
            fi
        else
            err_msg
        fi
    else
        err_msg
    fi
}

cp_file $1
