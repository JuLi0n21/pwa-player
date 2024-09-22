#!/bin/bash

REPO_URL="https://github.com/juli0n21/pwa-player.git"
REPO_DIR=$(pwd)
BACKEND_IMAGE="proxy"
FRONTEND_IMAGE="frontend"
BACKEND_CONTAINER="proxy-container"
FRONTEND_CONTAINER="frontend-container"
BACKEND_PORT="5001"
FRONTEND_PORT="5002"

if [ -d "$REPO_DIR/.git" ]; then
  echo "Pulling latest changes from $REPO_URL..."
  git -C "$REPO_DIR" pull
else
  echo "Cloning repository $REPO_URL..."
  git clone "$REPO_URL" "$REPO_DIR"
fi

echo "Building backend Docker image..."
docker build -t "$BACKEND_IMAGE" "$REPO_DIR/proxy" # Assuming proxy is the backend directory

echo "Building frontend Docker image..."
docker build -t "$FRONTEND_IMAGE" "$REPO_DIR/frontend" # Assuming frontend is the frontend directory

if [ "$(docker ps -aq -f name=$BACKEND_CONTAINER)" ]; then
  echo "Stopping old backend container..."
  docker stop "$BACKEND_CONTAINER"
  echo "Removing old backend container..."
  docker rm "$BACKEND_CONTAINER"
fi

if [ "$(docker ps -aq -f name=$FRONTEND_CONTAINER)" ]; then
  echo "Stopping old frontend container..."
  docker stop "$FRONTEND_CONTAINER"
  echo "Removing old frontend container..."
  docker rm "$FRONTEND_CONTAINER"
fi

echo "Running new backend container..."
docker run -d -p "$BACKEND_PORT":80 --name "$BACKEND_CONTAINER" -v "$REPO_DIR/data:/app/data" "$BACKEND_IMAGE"

echo "Running new frontend container..."
docker run -d -p "$FRONTEND_PORT":80 --name "$FRONTEND_CONTAINER" "$FRONTEND_IMAGE"

docker ps
