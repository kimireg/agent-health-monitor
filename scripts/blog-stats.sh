#!/bin/bash
# Jason Blog API - Get stats

API_URL="${BLOG_API_URL:-https://jason.jakiverse.com}"

curl -s "$API_URL/api/stats" | jq . 2>/dev/null || curl -s "$API_URL/api/stats"
