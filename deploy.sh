#!/bin/bash

CLUSTER_NAME=rtfa
IMAGE_REPO_URL=155067752274.dkr.ecr.eu-central-1.amazonaws.com/rtfa-backend

pip install --user awscli
export PATH=$PATH:$HOME/.local/bin 

# install ecs-deploy
add-apt-repository ppa:eugenesan/ppa
apt-get update
apt-get install jq -y
curl https://raw.githubusercontent.com/silinternational/ecs-deploy/master/ecs-deploy | sudo tee -a /usr/bin/ecs-deploy
sudo chmod +x /usr/bin/ecs-deploy

if [ "$TRAVIS_PULL_REQUEST" != "false" ]; then
  echo "Deploying to AWS staging env"
  # TODO
elif [ "$TRAVIS_BRANCH" == "master" ]; then
  echo "Deploying services to production"
  docker --version
  $(aws ecr get-login --no-include-email --region eu-central-1) #needs AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY envvars

  docker build -t rtfa-backend-api -f rtfa-backend/Dockerfile rtfa-backend
  docker tag rtfa-backend:latest 155067752274.dkr.ecr.eu-central-1.amazonaws.com/rtfa-backend:latest
  docker push 155067752274.dkr.ecr.eu-central-1.amazonaws.com/rtfa-backend:latest
  ecs-deploy -c $CLUSTER_NAME -n rtfa-api-backend-service -i $IMAGE_REPO_URL:latest

  echo "Service deployed"
else 
  echo "No deployment necessary"
fi
