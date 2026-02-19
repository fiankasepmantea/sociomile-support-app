#!/bin/bash
# reset-docker.sh
# Reset full Docker environment dan rebuild

echo "⚠️  Stopping and removing all containers, volumes, and orphaned networks..."
docker-compose down -v --remove-orphans

echo "🔧 Building images and starting containers..."
docker-compose up --build
