#!/bin/bash

TAG=$1
if [[ -z $TAG ]]
then
    TAG="latest"
fi

docker build -t feraru/crypto-bug:$TAG .
docker push feraru/crypto-bug:$TAG