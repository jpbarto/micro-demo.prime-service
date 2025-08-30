#!/bin/sh

docker run \
    -d --rm \
    --network=host \
    prime-service | tee results.log

sleep 3

curl -s 'http://localhost:8080/primes?number=10' | jq -c '.primes' | grep -E '^\[2,3,5,7\]$'
if [ $? -eq 0 ]; then
    echo "Test passed"
    docker stop $(docker ps -a -q --filter ancestor=prime-service:latest --format="{{.ID}}")
else
    echo "Primes service did not return expected result"
    docker stop $(docker ps -a -q --filter ancestor=prime-service:latest --format="{{.ID}}")
    exit 1
fi