#!/bin/bash

REPO_URL="https://github.com/juli0n21/pwa-player.git"
REPO_DIR=$(pwd)

if [ -d "$REPO_DIR/.git" ]; then
  echo "Pulling latest changes from $REPO_URL..."
  git -C "$REPO_DIR" pull
else
  echo "Cloning repository $REPO_URL..."
  git clone "$REPO_URL" "$REPO_DIR"
fi

docker compose up --build -d

