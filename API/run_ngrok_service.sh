#!/bin/bash
# Script to run ngrok as a permanent background service on EC2
# This will keep ngrok running even after you disconnect from the server

# Check if script is run with root privileges
if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run as root or with sudo"
    exit 1
fi

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
    echo "You can get your auth token by signing up at https://dashboard.ngrok.com/signup"
    exit 1
fi

# Create systemd service file for ngrok
cat > /tmp/ngrok.service << EOF
[Unit]
Description=ngrok tunnel service
After=network.target

[Service]
Type=simple
User=$(whoami)
WorkingDirectory=/home/$(whoami)
Environment="NGROK_AUTH_TOKEN=${NGROK_AUTH_TOKEN}"
ExecStart=/usr/bin/ngrok http --url=rational-evolving-joey.ngrok-free.app 8000 --log=stdout
Restart=always
RestartSec=10
# Keep the service running even if terminal closes
KillMode=process
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Move service file to systemd directory
sudo mv /tmp/ngrok.service /etc/systemd/system/

# Reload systemd, enable and start the service
sudo systemctl daemon-reload
sudo systemctl enable ngrok
sudo systemctl start ngrok

echo "ngrok service has been set up and started."
echo "To check the status, run: sudo systemctl status ngrok"
echo "To view the ngrok URL, run: curl http://localhost:4040/api/tunnels | jq '.tunnels[0].public_url'"
echo "Make sure to install jq with: sudo apt install -y jq"
