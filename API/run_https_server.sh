#!/bin/bash
# Script to run the FastAPI server with HTTPS

# Make sure the SSL directory exists
mkdir -p ssl

# Generate SSL certificates if they don't exist
if [ ! -f "ssl/cert.pem" ] || [ ! -f "ssl/key.pem" ]; then
    echo "Generating SSL certificates..."
    python generate_cert.py
fi

# Set environment variables
export USE_SSL=True
export PORT=8000

# Run the server
echo "Starting HTTPS server on port $PORT..."
python server.py
