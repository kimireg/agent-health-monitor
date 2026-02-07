# Agent Health Monitor

A lightweight web service for monitoring OpenClaw infrastructure status. Provides real-time visibility into email systems, Telegram gateway, cron jobs, and system resources.

## Features

- 📧 **Email System Monitoring** - Mirador running status, processor health, queue length, boot scan tracking
- 💬 **Telegram Gateway Status** - Connection status and health checks
- ⏰ **Cron Job Tracking** - List all scheduled tasks and last execution times
- 🖥️ **System Resources** - Disk and memory usage with visual progress bars
- 🎨 **Clean Dashboard** - Minimal, responsive web interface with auto-refresh
- 🔌 **REST API** - JSON endpoints for programmatic access

## API Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /` | HTML Dashboard |
| `GET /health` | Full health status (JSON) |
| `GET /health/email` | Email system details (JSON) |
| `GET /health/telegram` | Telegram status (JSON) |

## Local Development

### Prerequisites

- Go 1.21 or higher

### Run Locally

```bash
# Clone the repository
git clone git@github.com:kimireg/agent-health-monitor.git
cd agent-health-monitor

# Run the service
go run main.go

# Or build and run
go build -o agent-health-monitor
./agent-health-monitor
```

The service will start on port 8080. Access the dashboard at:
- Dashboard: http://localhost:8080
- Health API: http://localhost:8080/health

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `GO_ENV` | Environment mode | `development` |

## Deployment

### Zeabur (Recommended)

1. Fork this repository to your GitHub account
2. Log in to [Zeabur](https://zeabur.com)
3. Create a new project and deploy from GitHub
4. Select the `agent-health-monitor` repository
5. Zeabur will automatically detect the `zeabur.yaml` configuration

Alternatively, use the Zeabur CLI:

```bash
# Install Zeabur CLI
npm install -g zeabur

# Login
zeabur login

# Deploy
zeabur deploy
```

### Docker

```bash
# Build the image
docker build -t agent-health-monitor .

# Run the container
docker run -p 8080:8080 agent-health-monitor
```

### Manual Server Deployment

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o agent-health-monitor

# Copy to server
scp agent-health-monitor user@server:/opt/agent-health-monitor/
scp -r templates user@server:/opt/agent-health-monitor/

# Run on server
./agent-health-monitor
```

## Project Structure

```
agent-health-monitor/
├── main.go              # Go backend service
├── templates/
│   └── index.html       # Dashboard template
├── zeabur.yaml          # Zeabur deployment config
├── Dockerfile           # Docker build configuration
├── go.mod               # Go module definition
└── README.md            # This file
```

## Health Status Definitions

### Overall Status
- **healthy** - All systems operational
- **degraded** - Some non-critical issues
- **unhealthy** - Critical failures detected

### Component Status
- **healthy** - Operating normally
- **warning** - Performance degradation
- **error** - Service unavailable
- **critical** - Immediate attention required

## License

MIT License - See LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For issues and feature requests, please use the [GitHub Issues](https://github.com/kimireg/agent-health-monitor/issues) page.
