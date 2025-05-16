#!/bin/bash

# Your GitHub details
USERNAME="tsoodo"
REPO="git-switch"
VERSION="v1.0.0"

echo "Getting SHA256 hashes for gs binaries from $USERNAME/$REPO..."

# Download binaries
echo "Downloading binaries..."
curl -L -o gs-darwin-amd64 "https://github.com/$USERNAME/$REPO/releases/download/$VERSION/gs-darwin-amd64"
curl -L -o gs-darwin-arm64 "https://github.com/$USERNAME/$REPO/releases/download/$VERSION/gs-darwin-arm64"
curl -L -o gs-linux-amd64 "https://github.com/$USERNAME/$REPO/releases/download/$VERSION/gs-linux-amd64"
curl -L -o gs-linux-arm64 "https://github.com/$USERNAME/$REPO/releases/download/$VERSION/gs-linux-arm64"

# Get SHA256 hashes
echo -e "\nSHA256 hashes:"
echo "darwin-amd64: $(shasum -a 256 gs-darwin-amd64 | cut -d' ' -f1)"
echo "darwin-arm64: $(shasum -a 256 gs-darwin-arm64 | cut -d' ' -f1)"
echo "linux-amd64:  $(shasum -a 256 gs-linux-amd64 | cut -d' ' -f1)"
echo "linux-arm64:  $(shasum -a 256 gs-linux-arm64 | cut -d' ' -f1)"

# Clean up
rm gs-*

echo -e "\nUpdate your Formula/gs.rb with these hashes!"
