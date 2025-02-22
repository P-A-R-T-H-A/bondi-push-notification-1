#!/bin/bash

set -e # Exit immediately if a command exits with a non-zero status
set -u # Treat unset variables as an error

### Pull Branch
echo "Please enter branch name to pull:"
read -r branch
echo "you have entered branch:" "$branch"
sleep 1

cd ..

git stash; git branch -f origin/"$branch"; git checkout "$branch"; git pull;

cd -

echo "creating docker images"
docker build --no-cache -t bondi-push-notifier ../.

echo "killing the running docker"
docker system prune -f

# Check if there are any containers matching the pattern before attempting to kill them
if docker ps -a | grep -q 'bondi-api'; then
    docker ps -a | grep 'bondi-api' | awk '{print $1}' | xargs docker stop --time=5
    #docker ps -a | grep 'bondi-api' | awk '{print $1}' | xargs docker kill --time=5
    docker ps -a | grep 'bondi-api' | awk '{print $1}' | xargs docker rm
else
    echo "No containers matching the pattern 'bondi-api' found."
fi

echo "running the bondi-api using docker"
docker run -d --restart=unless-stopped --name bondi-push-notifier -p 8080:8080 bondi-push-notifier

echo "we are done!"