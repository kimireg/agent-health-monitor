#!/bin/bash
# Script to create GitHub repository using API
# Requires: GITHUB_TOKEN environment variable

if [ -z "$GITHUB_TOKEN" ]; then
    echo "Error: GITHUB_TOKEN environment variable not set"
    echo "Please set your GitHub personal access token:"
    echo "export GITHUB_TOKEN='your_token_here'"
    exit 1
fi

curl -X POST \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/user/repos \
  -d '{
    "name": "agent-health-monitor",
    "description": "A lightweight web service for monitoring OpenClaw infrastructure status",
    "private": false,
    "auto_init": false
  }'
