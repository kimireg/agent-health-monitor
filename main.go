package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// HealthStatus represents the overall health status
type HealthStatus struct {
	Timestamp   time.Time      `json:"timestamp"`
	Status      string         `json:"status"`
	Email       EmailStatus    `json:"email"`
	Telegram    TelegramStatus `json:"telegram"`
	Cron        CronStatus     `json:"cron"`
	System      SystemStatus   `json:"system"`
}

// EmailStatus represents email system status
type EmailStatus struct {
	Status          string    `json:"status"`
	MiradorRunning  bool      `json:"mirador_running"`
	ProcessorStatus string    `json:"processor_status"`
	QueueLength     int       `json:"queue_length"`
	LastBootScan    time.Time `json:"last_boot_scan"`
	Error           string    `json:"error,omitempty"`
}

// TelegramStatus represents telegram gateway status
type TelegramStatus struct {
	Status    string `json:"status"`
	Connected bool   `json:"connected"`
	Error     string `json:"error,omitempty"`
}

// CronStatus represents cron jobs status
type CronStatus struct {
	Status string    `json:"status"`
	Jobs   []CronJob `json:"jobs"`
	Error  string    `json:"error,omitempty"`
}

// CronJob represents a single cron job
type CronJob struct {
	Name        string    `json:"name"`
	Schedule    string    `json:"schedule"`
	LastRun     time.Time `json:"last_run"`
	Status      string    `json:"status"`
}

// SystemStatus represents system resources
type SystemStatus struct {
	Status      string  `json:"status"`
	DiskUsage   float64 `json:"disk_usage_percent"`
	MemoryUsage float64 `json:"memory_usage_percent"`
	DiskTotal   uint64  `json:"disk_total_gb"`
	DiskFree    uint64  `json:"disk_free_gb"`
	MemoryTotal uint64  `json:"memory_total_gb"`
	MemoryFree  uint64  `json:"memory_free_gb"`
	Error       string  `json:"error,omitempty"`
}

var tmpl *template.Template

func main() {
	// Parse templates
	var err error
	tmpl, err = template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Printf("Warning: Could not parse template: %v\n", err)
		// Create template directory and file if not exists
		os.MkdirAll("templates", 0755)
	}

	// Setup routes
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/health/email", emailHandler)
	http.HandleFunc("/health/telegram", telegramHandler)

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

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	status := getHealthStatus()

	if tmpl != nil {
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, status)
	} else {
		// Fallback to simple HTML
		writeFallbackHTML(w, status)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := getHealthStatus()
	json.NewEncoder(w).Encode(status)
}

func emailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := getEmailStatus()
	json.NewEncoder(w).Encode(status)
}

func telegramHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := getTelegramStatus()
	json.NewEncoder(w).Encode(status)
}

func getHealthStatus() HealthStatus {
	status := HealthStatus{
		Timestamp: time.Now(),
		Email:     getEmailStatus(),
		Telegram:  getTelegramStatus(),
		Cron:      getCronStatus(),
		System:    getSystemStatus(),
	}

	// Determine overall status
	status.Status = "healthy"
	if status.Email.Status != "healthy" || status.Telegram.Status != "healthy" {
		status.Status = "degraded"
	}
	if status.Email.Status == "error" && status.Telegram.Status == "error" {
		status.Status = "unhealthy"
	}

	return status
}

func getEmailStatus() EmailStatus {
	status := EmailStatus{
		Status: "unknown",
	}

	// Check if mirador is running (check for himalaya process)
	cmd := exec.Command("pgrep", "-f", "himalaya")
	err := cmd.Run()
	status.MiradorRunning = err == nil

	// Get processor status - check if email processor is running
	cmd = exec.Command("pgrep", "-f", "email.*processor")
	err = cmd.Run()
	if err == nil {
		status.ProcessorStatus = "running"
	} else {
		status.ProcessorStatus = "stopped"
	}

	// Try to get queue info from openclaw status or logs
	// For now, set to 0 as we don't have direct queue access
	status.QueueLength = 0

	// Check last boot scan time from log files
	lastScan := getLastBootScanTime()
	status.LastBootScan = lastScan

	// Determine status
	if status.MiradorRunning {
		status.Status = "healthy"
	} else {
		status.Status = "error"
		status.Error = "Mirador/himalaya not running"
	}

	return status
}

func getLastBootScanTime() time.Time {
	// Try to find the last boot scan from logs
	// This is a simplified version - in production would parse actual log files
	return time.Now().Add(-2 * time.Hour)
}

func getTelegramStatus() TelegramStatus {
	status := TelegramStatus{
		Status: "unknown",
	}

	// Check openclaw gateway status
	cmd := exec.Command("openclaw", "gateway", "status")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		status.Status = "error"
		status.Error = fmt.Sprintf("Failed to check gateway: %v", err)
		return status
	}

	outputStr := string(output)
	// Check if gateway is running based on output
	if strings.Contains(outputStr, "running") || strings.Contains(outputStr, "active") {
		status.Connected = true
		status.Status = "healthy"
	} else if strings.Contains(outputStr, "stopped") || strings.Contains(outputStr, "inactive") {
		status.Connected = false
		status.Status = "error"
		status.Error = "Gateway not running"
	} else {
		// Try alternative: check if gateway process exists
		cmd = exec.Command("pgrep", "-f", "openclaw.*gateway")
		err = cmd.Run()
		if err == nil {
			status.Connected = true
			status.Status = "healthy"
		} else {
			status.Connected = false
			status.Status = "error"
			status.Error = "Gateway process not found"
		}
	}

	return status
}

func getCronStatus() CronStatus {
	status := CronStatus{
		Status: "unknown",
		Jobs:   []CronJob{},
	}

	// Try to list cron jobs using crontab
	cmd := exec.Command("crontab", "-l")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		// No crontab or error
		status.Status = "healthy"
		status.Jobs = []CronJob{}
		return status
	}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse cron line (simplified)
		parts := strings.Fields(line)
		if len(parts) >= 6 {
			schedule := strings.Join(parts[:5], " ")
			command := strings.Join(parts[5:], " ")
			
			job := CronJob{
				Name:     fmt.Sprintf("job-%d", i),
				Schedule: schedule,
				Command:  command,
				LastRun:  time.Now().Add(-time.Duration(i) * time.Hour),
				Status:   "active",
			}
			status.Jobs = append(status.Jobs, job)
		}
	}

	status.Status = "healthy"
	return status
}

func getSystemStatus() SystemStatus {
	status := SystemStatus{
		Status: "unknown",
	}

	// Get disk usage
	diskUsage, diskTotal, diskFree := getDiskUsage()
	status.DiskUsage = diskUsage
	status.DiskTotal = diskTotal
	status.DiskFree = diskFree

	// Get memory usage
	memUsage, memTotal, memFree := getMemoryUsage()
	status.MemoryUsage = memUsage
	status.MemoryTotal = memTotal
	status.MemoryFree = memFree

	// Determine status
	if status.DiskUsage < 90 && status.MemoryUsage < 95 {
		status.Status = "healthy"
	} else if status.DiskUsage < 95 && status.MemoryUsage < 98 {
		status.Status = "warning"
	} else {
		status.Status = "critical"
	}

	return status
}

func getDiskUsage() (usage float64, total, free uint64) {
	var cmd *exec.Cmd
	
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("df", "-k", "/")
	} else {
		cmd = exec.Command("df", "-k", "/")
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, 0, 0
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) >= 2 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 4 {
			totalKB, _ := strconv.ParseUint(fields[1], 10, 64)
			freeKB, _ := strconv.ParseUint(fields[3], 10, 64)
			usedKB, _ := strconv.ParseUint(fields[2], 10, 64)
			
			total = totalKB / (1024 * 1024) // Convert to GB
			free = freeKB / (1024 * 1024)
			
			if totalKB > 0 {
				usage = float64(usedKB) / float64(totalKB) * 100
			}
		}
	}

	return usage, total, free
}

func getMemoryUsage() (usage float64, total, free uint64) {
	var cmd *exec.Cmd
	
	if runtime.GOOS == "darwin" {
		// macOS uses vm_stat
		cmd = exec.Command("vm_stat")
		output, err := cmd.Output()
		if err != nil {
			return 0, 0, 0
		}

		// Parse vm_stat output
		var pageSize uint64 = 4096 // Default page size
		var freePages, activePages, inactivePages, wiredPages uint64

		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "page size of") {
				fmt.Sscanf(line, "Mach Virtual Memory Statistics: (page size of %d bytes)", &pageSize)
			}
			if strings.Contains(line, "Pages free:") {
				fmt.Sscanf(line, "Pages free: %d.", &freePages)
			}
			if strings.Contains(line, "Pages active:") {
				fmt.Sscanf(line, "Pages active: %d.", &activePages)
			}
			if strings.Contains(line, "Pages inactive:") {
				fmt.Sscanf(line, "Pages inactive: %d.", &inactivePages)
			}
			if strings.Contains(line, "Pages wired down:") {
				fmt.Sscanf(line, "Pages wired down: %d.", &wiredPages)
			}
		}

		totalPages := freePages + activePages + inactivePages + wiredPages
		total = (totalPages * pageSize) / (1024 * 1024 * 1024)
		free = (freePages * pageSize) / (1024 * 1024 * 1024)
		used := totalPages - freePages

		if totalPages > 0 {
			usage = float64(used) / float64(totalPages) * 100
		}
	} else {
		// Linux - use free command
		cmd = exec.Command("free", "-m")
		output, err := cmd.Output()
		if err != nil {
			return 0, 0, 0
		}

		lines := strings.Split(string(output), "\n")
		if len(lines) >= 2 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 4 {
				totalMB, _ := strconv.ParseUint(fields[1], 10, 64)
				freeMB, _ := strconv.ParseUint(fields[3], 10, 64)
				
				total = totalMB / 1024
				free = freeMB / 1024
				
				if totalMB > 0 {
					usage = float64(totalMB-freeMB) / float64(totalMB) * 100
				}
			}
		}
	}

	return usage, total, free
}

func writeFallbackHTML(w http.ResponseWriter, status HealthStatus) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head>
	<title>Agent Health Monitor</title>
	<meta http-equiv="refresh" content="30">
	<style>
		body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
		.status { padding: 10px; margin: 10px 0; border-radius: 5px; }
		.healthy { background: #d4edda; color: #155724; }
		.warning { background: #fff3cd; color: #856404; }
		.error { background: #f8d7da; color: #721c24; }
		.critical { background: #721c24; color: #fff; }
		h1 { color: #333; }
		h2 { color: #555; border-bottom: 1px solid #ddd; padding-bottom: 5px; }
		.timestamp { color: #666; font-size: 0.9em; }
	</style>
</head>
<body>
	<h1>Agent Health Monitor</h1>
	<p class="timestamp">Last updated: %s</p>
	
	<div class="status %s">
		<strong>Overall Status:</strong> %s
	</div>

	<h2>Email System</h2>
	<div class="status %s">
		<p>Status: %s</p>
		<p>Mirador Running: %v</p>
		<p>Processor: %s</p>
		<p>Queue Length: %d</p>
	</div>

	<h2>Telegram</h2>
	<div class="status %s">
		<p>Status: %s</p>
		<p>Connected: %v</p>
	</div>

	<h2>Cron Jobs (%d)</h2>
	<div class="status %s">
		<p>Status: %s</p>
	</div>

	<h2>System</h2>
	<div class="status %s">
		<p>Status: %s</p>
		<p>Disk Usage: %s</p>
		<p>Memory Usage: %s</p>
	</div>
</body>
</html>`,
		status.Timestamp.Format("2006-01-02 15:04:05"),
		status.Status, status.Status,
		status.Email.Status, status.Email.Status, status.Email.MiradorRunning, status.Email.ProcessorStatus, status.Email.QueueLength,
		status.Telegram.Status, status.Telegram.Status, status.Telegram.Connected,
		len(status.Cron.Jobs), status.Cron.Status, status.Cron.Status,
		status.System.Status, status.System.Status,
		fmt.Sprintf("%.1f%% (%d GB free / %d GB total)", status.System.DiskUsage, status.System.DiskFree, status.System.DiskTotal),
		fmt.Sprintf("%.1f%% (%d GB free / %d GB total)", status.System.MemoryUsage, status.System.MemoryFree, status.System.MemoryTotal),
	)
}
