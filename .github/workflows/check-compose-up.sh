#!/bin/bash

echo "Waiting for application to be ready..."
timeout=5
counter=0
while ! curl -s http://localhost:8080/health > /dev/null; do
  if [ $counter -gt $timeout ]; then
    echo "Timeout waiting for application to be ready"
    docker compose logs
    exit 1
  fi
  echo "Still waiting..."
  sleep 1
  ((counter++))
done
echo "Application is ready!"