#!/bin/bash
# Script to start ngrok with the custom subdomain

# Check if NGROK_AUTH_TOKEN is set in .env file
if [ -f ".env" ]; then
    # Source the .env file to get environment variables
    source .env
fi

# Check if token is set
if [ -z "$NGROK_AUTH_TOKEN" ]; then
    echo "NGROK_AUTH_TOKEN environment variable not found."
    echo "Please add NGROK_AUTH_TOKEN=your_token to your .env file or export it in your shell."
    exit 1
fi

# Start ngrok with the custom subdomain
echo "Starting ngrok with custom subdomain: rational-evolving-joey.ngrok-free.app"
ngrok http --url=rational-evolving-joey.ngrok-free.app 8000
