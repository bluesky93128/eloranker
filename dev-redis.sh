#!/usr/bin/env bash

docker-machine start
echo "REDIS_ADDRESS=\"$(docker-machine ip):6379\""
eval $(docker-machine env)
docker-compose up redis
