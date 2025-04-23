#!/bin/bash

# Create data directory if it doesn't exist
mkdir -p data

# Build and start the containers
docker-compose up --build
