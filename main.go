package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"
	"time"
)

// PushStatus represents status data pushed from Mac Mini
type PushStatus struct {
	Timestamp string     `json:"timestamp"`
	Source    string     `json:"source"`
	Email     EmailPush  `json:"email"`
	Telegram  TgPush     `json:"telegram"`
	System    SysPush    `json:"system"`
}

type EmailPush struct {
	MiradorRunning  bool   `json:"mirador_running"`
	ProcessorRunning bool  `json:"processor_running"`
	UnreadCount     int    `json:"unread_count"`
}

type TgPush struct {
	GatewayRunning bool `json:"gateway_running"`
}

type SysPush struct {
	Uptime string `json:"uptime"`
}

// Global storage for pushed data
var (
	lastPush     *PushStatus
	lastPushTime time.Time
	pushMutex    sync.RWMutex
)

func main() {
	// Setup routes
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/api/health", healthHandler)
	http.HandleFunc("/api/push", pushHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Agent Health Monitor starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	}
}

func pushHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var status PushStatus
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	pushMutex.Lock()
	lastPush = &status
	lastPushTime = time.Now()
	pushMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "received"})
	fmt.Printf("[%s] Received push from %s\n", time.Now().Format("15:04:05"), status.Source)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	pushMutex.RLock()
	push := lastPush
	pushTime := lastPushTime
	pushMutex.RUnlock()

	response := map[string]interface{}{
		"timestamp":       time.Now().UTC(),
		"last_push":       pushTime,
		"last_push_data":  push,
		"status":          "waiting_for_data",
	}

	if push != nil {
		minutesSincePush := time.Since(pushTime).Minutes()
		if minutesSincePush < 15 {
			response["status"] = "healthy"
		} else {
			response["status"] = "stale_data"
		}
	}

	json.NewEncoder(w).Encode(response)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	pushMutex.RLock()
	push := lastPush
	pushTime := lastPushTime
	pushMutex.RUnlock()

	data := DashboardData{
		Timestamp:        time.Now(),
		LastPushTime:     pushTime,
		HasData:          push != nil,
		MinutesSincePush: int(time.Since(pushTime).Minutes()),
	}

	if push != nil {
		data.Source = push.Source
		data.MiradorRunning = push.Email.MiradorRunning
		data.ProcessorRunning = push.Email.ProcessorRunning
		data.UnreadCount = push.Email.UnreadCount
		data.GatewayRunning = push.Telegram.GatewayRunning
		data.Uptime = push.System.Uptime
		
		if data.MinutesSincePush < 15 {
			data.Status = "healthy"
			data.StatusColor = "#28a745"
		} else {
			data.Status = "stale"
			data.StatusColor = "#ffc107"
		}
	} else {
		data.Status = "waiting"
		data.StatusColor = "#6c757d"
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

// DashboardData for template
type DashboardData struct {
	Timestamp        time.Time
	LastPushTime     time.Time
	HasData          bool
	MinutesSincePush int
	Status           string
	StatusColor      string
	Source           string
	MiradorRunning   bool
	ProcessorRunning bool
	UnreadCount      int
	GatewayRunning   bool
	Uptime           string
}

var tmpl = template.Must(template.New("dashboard").Parse(`<!DOCTYPE html>
<html>
<head>
	<title>Agent Health Monitor</title>
	<meta http-equiv="refresh" content="30">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
		* { box-sizing: border-box; margin: 0; padding: 0; }
		body { 
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
			background: #f5f5f7;
			padding: 20px;
			color: #1d1d1f;
		}
		.container { max-width: 800px; margin: 0 auto; }
		header { text-align: center; margin-bottom: 30px; }
		h1 { font-size: 28px; margin-bottom: 10px; }
		.subtitle { color: #86868b; font-size: 14px; }
		.status-card {
			background: white;
			border-radius: 16px;
			padding: 24px;
			margin-bottom: 20px;
			box-shadow: 0 2px 8px rgba(0,0,0,0.08);
		}
		.status-header {
			display: flex;
			align-items: center;
			justify-content: space-between;
			margin-bottom: 16px;
		}
		.status-badge {
			display: inline-flex;
			align-items: center;
			gap: 8px;
			padding: 8px 16px;
			border-radius: 20px;
			font-size: 14px;
			font-weight: 500;
			color: white;
		}
		.status-dot {
			width: 8px;
			height: 8px;
			border-radius: 50%;
			background: white;
			animation: pulse 2s infinite;
		}
		@keyframes pulse {
			0%, 100% { opacity: 1; }
			50% { opacity: 0.5; }
		}
		.last-update {
			font-size: 13px;
			color: #86868b;
		}
		.service-grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
			gap: 16px;
		}
		.service-item {
			background: #f5f5f7;
			border-radius: 12px;
			padding: 16px;
		}
		.service-name {
			font-size: 12px;
			color: #86868b;
			text-transform: uppercase;
			letter-spacing: 0.5px;
			margin-bottom: 8px;
		}
		.service-value {
			font-size: 18px;
			font-weight: 600;
		}
		.service-value.running { color: #28a745; }
		.service-value.stopped { color: #dc3545; }
		.metric-value {
			font-size: 32px;
			font-weight: 700;
			color: #007aff;
		}
		.empty-state {
			text-align: center;
			padding: 60px 20px;
			color: #86868b;
		}
		.empty-state svg {
			width: 64px;
			height: 64px;
			margin-bottom: 16px;
			opacity: 0.5;
		}
		footer {
			text-align: center;
			margin-top: 40px;
			padding-top: 20px;
			border-top: 1px solid #d2d2d7;
			font-size: 12px;
			color: #86868b;
		}
	</style>
</head>
<body>
	<div class="container">
		<header>
			<h1>Agent Health Monitor</h1>
			<p class="subtitle">Jason's First Software</p>
		</header>

		{{if .HasData}}
		<div class="status-card">
			<div class="status-header">
				<span class="status-badge" style="background: {{.StatusColor}};">
					<span class="status-dot"></span>
					{{.Status}}
				</span>
				<span class="last-update">
					{{if eq .MinutesSincePush 0}}Just now{{else}}{{.MinutesSincePush}}m ago{{end}}
				</span>
			</div>
			
			<div class="service-grid">
				<div class="service-item">
					<div class="service-name">Mirador</div>
					<div class="service-value {{if .MiradorRunning}}running{{else}}stopped{{end}}">
						{{if .MiradorRunning}}[OK] Running{{else}}[NO] Stopped{{end}}
					</div>
				</div>
				<div class="service-item">
					<div class="service-name">Processor</div>
					<div class="service-value {{if .ProcessorRunning}}running{{else}}stopped{{end}}">
						{{if .ProcessorRunning}}[OK] Running{{else}}[NO] Stopped{{end}}
					</div>
				</div>
				<div class="service-item">
					<div class="service-name">Telegram Gateway</div>
					<div class="service-value {{if .GatewayRunning}}running{{else}}stopped{{end}}">
						{{if .GatewayRunning}}[OK] Connected{{else}}[NO] Disconnected{{end}}
					</div>
				</div>
				<div class="service-item">
					<div class="service-name">Unread Emails</div>
					<div class="metric-value">{{.UnreadCount}}</div>
				</div>
			</div>
			
			{{if .Uptime}}
			<div style="margin-top: 16px; padding-top: 16px; border-top: 1px solid #d2d2d7;">
				<span style="font-size: 13px; color: #86868b;">Mac Mini Uptime: {{.Uptime}}</span>
			</div>
			{{end}}
		</div>
		{{else}}
		<div class="status-card empty-state">
			<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
				<path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
			</svg>
			<h3>Waiting for data...</h3>
			<p>The Mac Mini hasn't reported status yet.<br>First push should arrive within 5 minutes.</p>
		</div>
		{{end}}

		<footer>
			<p>Last page refresh: {{.Timestamp.Format "15:04:05"}}</p>
			<p style="margin-top: 8px;">Jason 🍎 · OpenClaw · 2026</p>
		</footer>
	</div>
</body>
</html>`))
