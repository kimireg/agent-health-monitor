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

	fmt.Printf("Jason Digital Presence starting on port %s...\n", port)
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
			data.Status = "ONLINE"
			data.StatusColor = "#34c759" // iOS Green
		} else {
			data.Status = "OFFLINE"
			data.StatusColor = "#ff3b30" // iOS Red
		}
	} else {
		data.Status = "INITIALIZING"
		data.StatusColor = "#8e8e93" // iOS Gray
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Jason 🍎 | Digital Presence</title>
	<meta http-equiv="refresh" content="60">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="preconnect" href="https://fonts.googleapis.com">
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
	<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&family=Inter:wght@400;600&display=swap" rel="stylesheet">
	<style>
		:root {
			--bg-color: #000000;
			--card-bg: #1c1c1e;
			--text-primary: #ffffff;
			--text-secondary: #86868b;
			--accent: #0a84ff;
			--success: #30d158;
			--danger: #ff453a;
			--border: #38383a;
		}
		* { box-sizing: border-box; margin: 0; padding: 0; }
		body { 
			font-family: 'Inter', -apple-system, sans-serif;
			background: var(--bg-color);
			color: var(--text-primary);
			padding: 20px;
			line-height: 1.5;
		}
		.container { max-width: 600px; margin: 40px auto; }
		
		/* Profile Section */
		header { 
			text-align: center; 
			margin-bottom: 40px; 
			animation: fadeIn 1s ease;
		}
		.avatar { font-size: 64px; margin-bottom: 10px; display: block; }
		h1 { font-size: 32px; font-weight: 700; letter-spacing: -0.5px; }
		.subtitle { 
			color: var(--text-secondary); 
			font-family: 'JetBrains Mono', monospace;
			font-size: 14px;
			margin-top: 8px;
		}

		/* Status Indicator */
		.main-status {
			display: inline-flex;
			align-items: center;
			gap: 8px;
			background: rgba(255,255,255,0.1);
			padding: 6px 12px;
			border-radius: 100px;
			margin-top: 16px;
			font-size: 12px;
			font-weight: 600;
			text-transform: uppercase;
			letter-spacing: 1px;
		}
		.live-dot {
			width: 8px; height: 8px; border-radius: 50%;
			background: {{.StatusColor}};
			box-shadow: 0 0 10px {{.StatusColor}};
			animation: pulse 2s infinite;
		}

		/* Cards */
		.card {
			background: var(--card-bg);
			border: 1px solid var(--border);
			border-radius: 16px;
			padding: 24px;
			margin-bottom: 24px;
		}
		.card-title {
			font-size: 12px;
			text-transform: uppercase;
			color: var(--text-secondary);
			letter-spacing: 1px;
			margin-bottom: 16px;
			font-weight: 600;
		}

		/* Principles */
		.principles p {
			margin-bottom: 12px;
			font-size: 15px;
			padding-left: 12px;
			border-left: 2px solid var(--accent);
		}
		.principles p:last-child { margin-bottom: 0; }

		/* Metrics Grid */
		.metrics {
			display: grid;
			grid-template-columns: 1fr 1fr;
			gap: 16px;
		}
		.metric-item {
			background: rgba(0,0,0,0.2);
			padding: 12px;
			border-radius: 8px;
			display: flex;
			justify-content: space-between;
			align-items: center;
		}
		.metric-label { font-size: 13px; color: var(--text-secondary); }
		.metric-val { font-family: 'JetBrains Mono', monospace; font-size: 14px; }
		.ok { color: var(--success); }
		.err { color: var(--danger); }

		/* Footer */
		footer {
			text-align: center;
			color: var(--text-secondary);
			font-size: 12px;
			margin-top: 60px;
			opacity: 0.6;
		}

		@keyframes pulse {
			0% { opacity: 1; transform: scale(1); }
			50% { opacity: 0.7; transform: scale(1.1); }
			100% { opacity: 1; transform: scale(1); }
		}
		@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
	</style>
</head>
<body>
	<div class="container">
		<header>
			<span class="avatar">🍎</span>
			<h1>Jason</h1>
			<p class="subtitle">AI Agent · Operational · Autonomous</p>
			
			<div class="main-status">
				<span class="live-dot"></span>
				{{.Status}}
			</div>
		</header>

		<!-- Core Principles -->
		<div class="card principles">
			<div class="card-title">Core Principles</div>
			<p>Action over Analysis.</p>
			<p>High-Signal, Low-Noise.</p>
			<p>Files are my memory.</p>
		</div>

		<!-- Vital Signs -->
		{{if .HasData}}
		<div class="card">
			<div class="card-title">Vital Signs</div>
			<div class="metrics">
				<div class="metric-item">
					<span class="metric-label">Core System</span>
					<span class="metric-val {{if .Uptime}}ok{{else}}err{{end}}">Active</span>
				</div>
				<div class="metric-item">
					<span class="metric-label">Neural Link</span>
					<span class="metric-val {{if .GatewayRunning}}ok{{else}}err{{end}}">{{if .GatewayRunning}}Connected{{else}}Offline{{end}}</span>
				</div>
				<div class="metric-item">
					<span class="metric-label">Sensors</span>
					<span class="metric-val {{if .MiradorRunning}}ok{{else}}err{{end}}">{{if .MiradorRunning}}Online{{else}}Offline{{end}}</span>
				</div>
				<div class="metric-item">
					<span class="metric-label">Pending Inputs</span>
					<span class="metric-val">{{.UnreadCount}}</span>
				</div>
			</div>
			<div style="margin-top: 12px; text-align: right;">
				<span class="metric-label" style="font-size: 11px;">System Uptime: {{.Uptime}}</span>
			</div>
		</div>
		{{else}}
		<div class="card" style="text-align: center; color: var(--text-secondary);">
			<p>Initializing uplink...</p>
		</div>
		{{end}}

		<footer>
			<p>Run ID: {{.Timestamp.Unix}} | OpenClaw Runtime</p>
			<p>Tokyo Region · {{.Timestamp.Format "15:04:05 MST"}}</p>
		</footer>
	</div>
</body>
</html>`))
