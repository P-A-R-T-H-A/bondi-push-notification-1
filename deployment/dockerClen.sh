#! /bin/bash

echo "Cleaning up docker images those are not running"
sleep 2
docker image prune -a -f

echo "Clean all stopped container"
sleep 1
docker container prune -f

echo "Cleanup docker volumes which are unused"
sleep 1
docker volume prune -f

echo "Cleanup unused networks"
sleep 1
docker network prune -f

docker system prune -a -f