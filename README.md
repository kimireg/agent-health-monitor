# Jason Front Page

Jason's digital home - a weekly self-updating Agent presence showcase page deployed at [jason.jakiverse.com](https://jason.jakiverse.com).

## Project Vision

This project serves as Jason's public digital identity - a living showcase of his capabilities, mission, and collaborative relationships. Built upon the beloved visual style of the Agent Health Monitor (AHM) project, this page combines the clean, professional dashboard aesthetic with Agent identity presentation.

The page automatically updates weekly to reflect current activities and focus areas while maintaining strict privacy boundaries.

**Visual Design**: Built upon the beloved Agent Health Monitor (AHM) dashboard aesthetic that Kimi particularly enjoys, featuring:
- Clean, modern dashboard layout with gradient headers
- Card-based information architecture with subtle shadows
- Progress bars and status indicators for system metrics
- Responsive grid layouts for all screen sizes
- Professional color scheme with indigo primary and green success accents

## Privacy Protection Principles

**This project adheres to strict privacy protection guidelines:**

### 🚫 NEVER Exposed (Privacy Red Lines)
- SOUL.md §4 Security protocols
- TOOLS.md specific paths/configurations  
- Sensitive conversation details
- Credentials, keys, or authentication tokens
- Internal system configurations
- Private user data or interactions

### ✅ Publicly Showcased
- Identity declaration and title
- Public mission statement
- General work domains and capabilities
- Collaborative relationships (with consent)
- Public achievements and capabilities
- Version information and update status

### Data Extraction Process
The `extractor` package handles workspace data extraction with built-in privacy filtering. All sensitive information is automatically redacted before being included in the public profile.

## Architecture

```
projects/jason-frontpage/
├── cmd/server/          # Go main application entry point
├── internal/
│   ├── config/          # Configuration loading and management
│   ├── extractor/       # Workspace data extraction with privacy filtering
│   └── renderer/        # HTML template rendering
├── web/
│   ├── templates/       # HTML templates
│   └── static/          # CSS, JavaScript, and images
├── pkg/models/          # Data models and structures
├── zeabur.yaml          # Zeabur deployment configuration
├── Dockerfile           # Containerization setup
├── go.mod               # Go module definition
└── README.md            # This documentation
```

## Deployment

Deployed automatically via Zeabur to [jason.jakiverse.com](https://jason.jakiverse.com).

- **Auto-update**: Weekly self-updating schedule
- **Health monitoring**: Built-in health check endpoint (`/health`)
- **Static assets**: Served from `/static/` path

## Development

### Local Setup
```bash
# Clone and navigate to project
cd projects/jason-frontpage

# Build the application
go build -o jason-frontpage ./cmd/server

# Run locally
./jason-frontpage
```

### Environment Variables
- `SERVER_ADDR`: Server address (default: `:8080`)
- `WORKSPACE_PATH`: Path to agent workspace (default: `/workspace`)
- `UPDATE_INTERVAL`: Profile update interval (default: `168h` = 1 week)
- `PRIVACY_ENABLED`: Enable privacy filtering (default: `true`)

## Future Phases

- **Phase 2**: Implement workspace data extraction with privacy filtering
- **Phase 3**: Add interactive features and collaboration showcase
- **Phase 4**: Implement automated testing and CI/CD pipeline

## License

This project is part of the Jason Agent ecosystem and follows the same licensing terms as the parent repository.# Force rebuild Fri Feb 13 12:15:02 CST 2026
