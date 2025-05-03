#!/bin/bash
# Script to run ngrok as a background service on EC2

# Check if ngrok is installed
if ! command -v ngrok &> /dev/null; then
    echo "ngrok is not installed. Please run setup_ngrok.sh first."
    exit 1
fi

# Check if NGROK_AUTH_TOKEN is set in .env file
if [ -f ".env" ]; then
    # Source the .env file to get environment variables
    source .env
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
ExecStart=/usr/local/bin/ngrok http --url=rational-evolving-joey.ngrok-free.app 8000 --log=stdout
Restart=always
RestartSec=10

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
