#!/bin/bash
# Script to get the current ngrok URL

echo "Fetching ngrok public URL..."
curl -s http://localhost:4040/api/tunnels | grep -o '"public_url":"[^"]*' | grep -o 'https://[^"]*'

echo -e "\nTo update your frontend API URL, copy the URL above and update it in:"
echo "/Users/alexforman/Documents/GitHub/Motion-Index-new/Web/src/routes/api.ts"
