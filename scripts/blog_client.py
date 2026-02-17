#!/usr/bin/env python3
"""
Jason Blog API Client
Simple client for posting to Jason's blog via API
"""

import os
import json
import urllib.request
import urllib.error
from datetime import datetime

class BlogClient:
    def __init__(self, api_url=None, api_key=None):
        self.api_url = api_url or os.getenv('BLOG_API_URL', 'https://jason.jakiverse.com')
        self.api_key = api_key or os.getenv('BLOG_API_KEY', '')
    
    def create_post(self, title, content, tags="daily", mood="focused", date=None):
        """Create a new blog post"""
        if date is None:
            date = datetime.now().strftime('%Y-%m-%d')
        
        data = {
            "date": date,
            "title": title,
            "content": content,
            "tags": tags,
            "mood": mood
        }
        
        url = f"{self.api_url}/api/posts"
        headers = {
            "Content-Type": "application/json"
        }
        if self.api_key:
            headers["X-API-Key"] = self.api_key
        
        req = urllib.request.Request(
            url,
            data=json.dumps(data).encode('utf-8'),
            headers=headers,
            method='POST'
        )
        
        try:
            with urllib.request.urlopen(req, timeout=30) as resp:
                return json.loads(resp.read().decode('utf-8'))
        except urllib.error.HTTPError as e:
            return {"error": e.read().decode('utf-8')}
        except Exception as e:
            return {"error": str(e)}
    
    def get_post(self, date=None):
        """Get a post by date (defaults to today)"""
        if date is None:
            date = datetime.now().strftime('%Y-%m-%d')
        
        url = f"{self.api_url}/api/post?date={date}"
        
        try:
            with urllib.request.urlopen(url, timeout=30) as resp:
                return json.loads(resp.read().decode('utf-8'))
        except urllib.error.HTTPError as e:
            if e.code == 404:
                return None
            return {"error": e.read().decode('utf-8')}
        except Exception as e:
            return {"error": str(e)}
    
    def list_posts(self, limit=100):
        """List all posts"""
        url = f"{self.api_url}/api/posts"
        
        try:
            with urllib.request.urlopen(url, timeout=30) as resp:
                return json.loads(resp.read().decode('utf-8'))
        except Exception as e:
            return {"error": str(e)}
    
    def get_stats(self):
        """Get blog statistics"""
        url = f"{self.api_url}/api/stats"
        
        try:
            with urllib.request.urlopen(url, timeout=30) as resp:
                return json.loads(resp.read().decode('utf-8'))
        except Exception as e:
            return {"error": str(e)}


# Simple usage example
if __name__ == "__main__":
    import sys
    
    if len(sys.argv) < 3:
        print("Usage: python blog_client.py \"Title\" \"Content\" [tags] [mood]")
        sys.exit(1)
    
    client = BlogClient()
    result = client.create_post(
        title=sys.argv[1],
        content=sys.argv[2],
        tags=sys.argv[3] if len(sys.argv) > 3 else "daily",
        mood=sys.argv[4] if len(sys.argv) > 4 else "focused"
    )
    
    print(json.dumps(result, indent=2))
