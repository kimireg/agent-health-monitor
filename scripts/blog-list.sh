#!/bin/bash
# Jason Blog API - List all posts

API_URL="${BLOG_API_URL:-https://jason.jakiverse.com}"

curl -s "$API_URL/api/posts" | jq . 2>/dev/null || curl -s "$API_URL/api/posts"
