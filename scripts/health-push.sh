#!/bin/bash
# Agent Health Monitor - Local Push Script
# Run on Mac Mini every 5 minutes via cron

ZEBUR_URL="https://ahm.zeabur.app/api/push"
# Simple auth token (can be enhanced later)
AUTH_TOKEN="jason-push-2026"

# Check services
MIRADOR_RUNNING=false
PROCESSOR_RUNNING=false
GATEWAY_RUNNING=false
UNREAD_COUNT=0

# Check Mirador
if pgrep -f "mirador watch" > /dev/null 2>&1; then
    MIRADOR_RUNNING=true
fi

# Check Processor
if pgrep -f "processor-daemon" > /dev/null 2>&1; then
    PROCESSOR_RUNNING=true
fi

# Check Telegram Gateway
if pgrep -f "openclaw.*gateway" > /dev/null 2>&1; then
    GATEWAY_RUNNING=true
fi

# Get unread email count (best effort)
if command -v himalaya &> /dev/null; then
    UNREAD_COUNT=$(himalaya envelope list --page-size 100 --output json 2>/dev/null | tail -1 | jq '[.[] | select(.flags | contains(["Seen"]) | not)] | length' 2>/dev/null || echo "0")
fi

# Get uptime
UPTIME=$(uptime | awk '{print $3}' | sed 's/,//')

# Build JSON payload
JSON=$(cat <<EOF
{
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "source": "mac-mini",
  "email": {
    "mirador_running": $MIRADOR_RUNNING,
    "processor_running": $PROCESSOR_RUNNING,
    "unread_count": $UNREAD_COUNT
  },
  "telegram": {
    "gateway_running": $GATEWAY_RUNNING
  },
  "system": {
    "uptime": "$UPTIME"
  }
}
EOF
)

# Push to Zeabur
RESPONSE=$(curl -s -X POST "$ZEBUR_URL" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -d "$JSON" 2>&1)

# Log locally (optional)
echo "[$(date '+%Y-%m-%d %H:%M:%S')] Push result: $RESPONSE"
