#!/bin/bash
docker-compose -f ./broker/docker-compose.yml up -d
sleep 10
docker run -v "$(pwd)":/go/github.com/petergood/codecollabbackend --env BOOTSTRAP_URLS=172.18.0.11:9092 --network broker_test ccb-broker-test
docker-compose -f ./broker/docker-compose.yml down