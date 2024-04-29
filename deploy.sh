#!/bin/bash

echo What should the version be?
read VERSION

docker build -t shariqalidev/discovery-trail:$VERSION .
docker push shariqalidev/discovery-trail:$VERSION