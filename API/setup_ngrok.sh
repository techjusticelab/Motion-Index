#!/bin/bash
# Script to set up ngrok on EC2 instance

# Install ngrok
echo "Installing ngrok..."
curl -s https://ngrok-agent.s3.amazonaws.com/ngrok.asc | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null
echo "deb https://ngrok-agent.s3.amazonaws.com buster main" | sudo tee /etc/apt/sources.list.d/ngrok.list
sudo apt update
sudo apt install -y ngrok

# Create ngrok config directory
mkdir -p ~/.ngrok2

# Check if NGROK_AUTH_TOKEN is set in .env file
if [ -f ".env" ]; then
    # Source the .env file to get environment variables
    source .env
fi

# Create ngrok configuration file
cat > ~/.ngrok2/ngrok.yml << EOF
version: "2"
authtoken: ${NGROK_AUTH_TOKEN:-YOUR_NGROK_AUTH_TOKEN}
tunnels:
  motion_index_api:
    proto: http
    addr: 8000
    subdomain: motion-index-api
EOF

# Check if token is set
if [ -z "$NGROK_AUTH_TOKEN" ]; then
    echo "NGROK_AUTH_TOKEN environment variable not found."
    echo "Please add NGROK_AUTH_TOKEN=your_token to your .env file or export it in your shell."
    echo "You can get your auth token by signing up at https://dashboard.ngrok.com/signup"
else
    echo "ngrok configured with auth token from environment variable."
fi
echo ""
echo "To start ngrok with a persistent tunnel, run:"
echo "ngrok start --all"
echo ""
echo "For a free account with a random URL, run:"
echo "ngrok http 8000"
