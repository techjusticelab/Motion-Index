#!/bin/bash
# Script to run ngrok as a persistent background process using nohup
# This will keep ngrok running even after you disconnect from the server

# Check if ngrok is installed
if ! command -v ngrok &> /dev/null; then
    echo "ngrok is not installed. Please run setup_ngrok.sh first."
    exit 1
fi

# Get the absolute path to the .env file
ENV_FILE="$(pwd)/.env"

# Check if NGROK_AUTH_TOKEN is set in .env file
if [ -f "$ENV_FILE" ]; then
    # Source the .env file to get environment variables
    source "$ENV_FILE"
    echo "Loaded environment variables from $ENV_FILE"
fi

# Check if token is set
if [ -z "$NGROK_AUTH_TOKEN" ]; then
    echo "NGROK_AUTH_TOKEN environment variable not found."
    echo "Please add NGROK_AUTH_TOKEN=your_token to your .env file or export it in your shell."
    exit 1
fi

# Create logs directory if it doesn't exist
mkdir -p logs

# Kill any existing ngrok processes
pkill -f ngrok || true
echo "Killed any existing ngrok processes"

# Start ngrok in the background with nohup
echo "Starting ngrok in the background with nohup..."
nohup ngrok http --url=rational-evolving-joey.ngrok-free.app 8000 > logs/ngrok.log 2>&1 &

# Save the process ID
echo $! > logs/ngrok.pid
echo "ngrok started with PID $(cat logs/ngrok.pid)"
echo "Logs are being written to logs/ngrok.log"
echo "To check if ngrok is running: ps -p $(cat logs/ngrok.pid)"
echo "To stop ngrok: kill $(cat logs/ngrok.pid)"

# Wait a moment for ngrok to start
sleep 3

# Check if ngrok is running
if ps -p $(cat logs/ngrok.pid) > /dev/null; then
    echo "ngrok is running successfully!"
    echo "You can now close this terminal, and ngrok will continue running."
    echo "To get the ngrok URL: ./get_ngrok_url.sh"
else
    echo "Failed to start ngrok. Check logs/ngrok.log for details."
fi
