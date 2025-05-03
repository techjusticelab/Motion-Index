#!/bin/bash
# Script to run the FastAPI server (HTTP only, HTTPS handled by ngrok)

# Set environment variables
export PORT=8000

# Run the server
echo "Starting HTTP server on port $PORT (HTTPS handled by ngrok)..."
python3 server.py
