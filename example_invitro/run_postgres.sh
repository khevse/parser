#!/bin/bash

WORKING_DIR=$(pwd)/pg_data

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
  -p5432:5432 \
  -v "$WORKING_DIR:/var/pgdata" \
  -e PGDATA=/var/pgdata \
  -e POSTGRES_INITDB_ARGS="--data-checksums --encoding=UTF-8" \
 postgres:9.6.2