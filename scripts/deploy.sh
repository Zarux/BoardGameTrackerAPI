#!/usr/bin/env bash

NAME=bgserver
IMAGE=${DOCKER_USERNAME}/${NAME}
TAG=""
PORT=8080

echo "Running deploy script"

if [[ $1 = "travis" ]]; then
    if [[ $2 = "dev" ]]; then
        echo "Using dev build"
        TAG=dev
        PORT=8081
    fi

    if [[ $2 = "prod" ]]; then
        echo "Using prod build"
        TAG=latest
    fi

    if [[ !  -z  ${TAG}  ]]; then
        echo "Building ${IMAGE}:${TAG}"
        docker build --no-cache -t ${IMAGE}:${TAG} . || (echo "Could not build" && exit 1)
        echo "Pushing ${IMAGE}:${TAG}"
        docker push ${IMAGE}:${TAG} || (echo "Could not push" && exit 1)
        echo "Deploying ${IMAGE}:${TAG}"
        sshpass -p ${SSH_PASS} ssh -o StrictHostKeyChecking=no ${SSH_USER}@${SSH_HOST} "docker stop ${NAME}-${TAG};docker rm -f ${NAME}-${TAG};docker pull ${IMAGE}:${TAG};docker run -d -p ${PORT}:8080 --restart always --name ${NAME}-${TAG} -e DSN_MAIN=${DSN_MAIN} ${IMAGE}:${TAG}" || (echo "Could not deploy" && exit 1)
    fi
fi
