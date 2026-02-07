# GitHub Repository Setup

The repository needs to be created on GitHub before pushing.

## Option 1: Web Interface (Easiest)

1. Go to https://github.com/new
2. Enter repository name: `agent-health-monitor`
3. Set visibility: Public
4. Click "Create repository"
5. Then run: `git push -u origin main`

## Option 2: GitHub CLI

```bash
# Install gh if not already installed
brew install gh

# Authenticate
gh auth login

# Create repository
gh repo create agent-health-monitor --public --source=. --push
```

## Option 3: GitHub API (with token)

```bash
# Set your GitHub token
export GITHUB_TOKEN='your_personal_access_token_here'

# Run the helper script
./create-repo.sh

# Then push
git push -u origin main
```

## After Repository Creation

Once the repository exists on GitHub, push the code:

```bash
cd /Users/kimi/.openclaw/workspace/research/lab/agent-health-monitor
git push -u origin main
```
