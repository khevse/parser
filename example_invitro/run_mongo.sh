#!/bin/bash

WORKING_DIR=$(pwd)/mongo_data

if [ -d "$WORKING_DIR" ]; then
    read -r -p "Directory '$WORKING_DIR' already exist. Remove? [Y/n]" response
    response=${response,,} # tolower

    if [[ $response =~ ^(yes|y)$ ]]; then
        sudo rm -Rf "$WORKING_DIR";
    fi
fi

if [ ! -d "$WORKING_DIR" ]; then
    mkdir "$WORKING_DIR"
fi

docker run -it --rm \
  -v "$WORKING_DIR:/data/db" \
  -p 27017:27017 \
mongo:3.4.2