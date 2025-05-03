#!/bin/bash
# Script to clean up disk space on EC2 instance

echo "Checking current disk space..."
df -h

echo -e "\nLargest directories in /var:"
du -h --max-depth=1 /var | sort -hr | head -10

echo -e "\nLargest directories in /tmp:"
du -h --max-depth=1 /tmp | sort -hr | head -10

echo -e "\nCleaning Docker resources..."
docker system prune -af --volumes

echo -e "\nCleaning package manager cache..."
apt-get clean
apt-get autoremove -y

echo -e "\nRemoving old log files..."
find /var/log -type f -name "*.gz" -delete
find /var/log -type f -name "*.log.*" -delete

echo -e "\nRemoving temporary files..."
rm -rf /tmp/*

echo -e "\nDisk space after cleanup:"
df -h
