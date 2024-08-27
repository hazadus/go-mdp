#!/bin/bash

./bin/mdp -file $1
FILEHASH=`md5 -r $1`
while true; do
    NEWHASH=`md5 -r $1`
    if [ "$NEWHASH" != "$FILEHASH" ]; then
        ./bin/mdp -file $1 > /dev/null
        FILEHASH=$NEWHASH
        echo "Reloading $1..."
    fi
    sleep 5
done
