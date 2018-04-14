#/bin/bash

image=rerorero/myghome
tag=$(git rev-parse HEAD)

set -eux
docker build -t $image:$tag ./
docker tag $image:$tag $image:latest
# push
docker login -u rerorero -p "$DOCKER_PASSWORD";
docker push $image:$tag
docker push $image:latest
