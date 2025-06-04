#!/bin/bash

for peer in $(docker ps --format "{{.Names}}" --filter "name=peer"); do
    echo "Installing Go inside $peer..."
    docker exec -u root -it $peer sh -c "
        if [ -f /etc/alpine-release ]; then
            apk add --no-cache go;
        elif [ -f /etc/redhat-release ]; then
            yum install -y golang;
        elif [ -f /etc/debian_version ]; then
            apt update && apt install -y golang;
        else
            echo 'Unsupported OS';
            exit 1;
        fi
        mkdir -p /go/pkg/mod &&
        chmod -R 777 /go /go/pkg /go/pkg/mod"
    docker restart $peer
done

