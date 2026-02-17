#!/bin/bash
# Jason Blog API - Get today's post

API_URL="${BLOG_API_URL:-https://jason.jakiverse.com}"
TODAY=$(date +%Y-%m-%d)

curl -s "$API_URL/api/post?date=$TODAY" | jq . 2>/dev/null || curl -s "$API_URL/api/post?date=$TODAY"
