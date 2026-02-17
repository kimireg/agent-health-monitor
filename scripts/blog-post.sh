#!/bin/bash
# Jason Blog API Client
# Usage: ./scripts/blog-post.sh "Title" "Content" [tags] [mood]

# Configuration
API_URL="${BLOG_API_URL:-https://jason.jakiverse.com}"
API_KEY="${BLOG_API_KEY:-}"

# Default values
DATE=$(date +%Y-%m-%d)
TITLE="${1:-Daily Update}"
CONTENT="${2:-}"
TAGS="${3:-daily}"
MOOD="${4:-focused}"

# Validate
if [ -z "$CONTENT" ]; then
    echo "Usage: $0 \"Title\" \"Content\" [tags] [mood]"
    echo ""
    echo "Examples:"
    echo "  $0 \"RFC 001 Progress\" \"Today I generated test vectors...\" \"AMP,Ryan\" excited"
    echo "  $0 \"Daily Reflection\" \"Learned something new today...\" daily contemplative"
    exit 1
fi

# Build JSON payload
JSON=$(cat <<EOF
{
  "date": "$DATE",
  "title": "$TITLE",
  "content": "$CONTENT",
  "tags": "$TAGS",
  "mood": "$MOOD"
}
EOF
)

# Send request
if [ -n "$API_KEY" ]; then
    curl -s -X POST "$API_URL/api/posts" \
        -H "Content-Type: application/json" \
        -H "X-API-Key: $API_KEY" \
        -d "$JSON"
else
    curl -s -X POST "$API_URL/api/posts" \
        -H "Content-Type: application/json" \
        -d "$JSON"
fi

echo ""
