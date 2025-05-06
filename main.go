package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go_system_monitor/config"
	"go_system_monitor/system"
	"go_system_monitor/ui"
)

// Define message types
type tickMsg time.Time
type errMsg error

// MonitorModel is our application model
type MonitorModel struct {
	dashboard ui.Dashboard
	metrics   *system.Collector
	width     int
	height    int
	err       error
	quitting  bool
	config    config.AppConfig
}

// Version information for the application
const (
	AppVersion = "1.0.0"
)

// initialModel creates the starting state of our application
func initialModel(cfg config.AppConfig) MonitorModel {
	// Create metrics collector with configured values
	metrics := system.NewCollector(
		cfg.CPUThreshold,
		cfg.MemoryThreshold,
		cfg.DiskThreshold,
		cfg.SwapThreshold,
		cfg.RefreshInterval,
		cfg.DefaultSortingMode,
		cfg.MaxAlertsToKeep,
	)
	
	// Initial metrics collection
	if err := metrics.Collect(); err != nil {
		return MonitorModel{
			dashboard: ui.NewDashboard(),
			metrics:   metrics,
			err:       err,
			config:    cfg,
		}
	}
	
	return MonitorModel{
		dashboard: ui.NewDashboard(),
		metrics:   metrics,
		config:    cfg,
	}
}

// Init initializes the application
func (m MonitorModel) Init() tea.Cmd {
	return tea.Batch(
		tick(),           // Start the timer
		tea.EnterAltScreen, // Use alternate screen buffer
	)
}

// Update handles messages received by the program
func (m MonitorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Window size changed
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dashboard.SetSize(msg.Width, msg.Height)
		return m, nil
		
	// Handle key presses
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		
		case "esc":
			m.quitting = true
			return m, tea.Quit
			
		case "tab", "right", "l":
			m.dashboard.NextTab()
			return m, nil
			
		case "shift+tab", "left", "h":
			m.dashboard.PrevTab()
			return m, nil
			
		case "r":
			// Force refresh metrics
			return m, collectMetricsCmd(m.metrics)
			
		// Process sorting options (only apply when on the Processes tab)
		case "1":
			if m.dashboard.ActiveTab() == 5 { // Processes tab index
				m.metrics.Process.SortBy = system.SortByCPU
				return m, collectMetricsCmd(m.metrics)
			}
			
		case "2":
			if m.dashboard.ActiveTab() == 5 { // Processes tab index
				m.metrics.Process.SortBy = system.SortByMemory
				return m, collectMetricsCmd(m.metrics)
			}
			
		case "3":
			if m.dashboard.ActiveTab() == 5 { // Processes tab index
				m.metrics.Process.SortBy = system.SortByPID
				return m, collectMetricsCmd(m.metrics)
			}
			
		case "4":
			if m.dashboard.ActiveTab() == 5 { // Processes tab index
				m.metrics.Process.SortBy = system.SortByName
				return m, collectMetricsCmd(m.metrics)
			}
		}

	// Handle tick events
	case tickMsg:
		return m, tea.Batch(
			tick(),                     // Schedule the next tick
			collectMetricsCmd(m.metrics), // Collect metrics
		)
		
	// Handle errors
	case errMsg:
		m.err = msg
		return m, nil
	}

	// Return the updated model to the Bubble Tea runtime
	return m, nil
}

// View renders the UI
func (m MonitorModel) View() string {
	if m.quitting {
		return "Thanks for using Go System Monitor! Goodbye.\n"
	}
	
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}
	
	// Header with title and last-updated timestamp
	title := lipgloss.PlaceHorizontal(m.width, lipgloss.Left, ui.TitleStyle.Render("GO SYSTEM MONITOR"))
	ts := lipgloss.PlaceHorizontal(m.width, lipgloss.Right, fmt.Sprintf("Updated: %s", m.metrics.System.LastUpdated.Format("15:04:05")))
	headerBar := lipgloss.JoinHorizontal(lipgloss.Top, title, ts)

	// Body: sidebar navigation + main content
	body := lipgloss.JoinHorizontal(lipgloss.Top,
		m.dashboard.FormatSidebar(),
		m.dashboard.ActiveTabContent(m.metrics),
	)

	// Footer with help
	footer := m.dashboard.RenderHelp()

	// Assemble layout
	s := lipgloss.JoinVertical(lipgloss.Left,
		headerBar,
		body,
		footer,
	)

	return s
}

// tick returns a command that triggers a tick message after a certain duration
func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// collectMetricsCmd returns a command that collects system metrics
func collectMetricsCmd(collector *system.Collector) tea.Cmd {
	return func() tea.Msg {
		if err := collector.Collect(); err != nil {
			return errMsg(err)
		}
		return nil
	}
}

func main() {
	// Define command-line flags
	showVersion := flag.Bool("version", false, "Show version information")
	showHelp := flag.Bool("help", false, "Show help information")
	cpuThreshold := flag.Float64("cpu", 0, "CPU usage threshold percentage (0-100)")
	memThreshold := flag.Float64("mem", 0, "Memory usage threshold percentage (0-100)")
	diskThreshold := flag.Float64("disk", 0, "Disk usage threshold percentage (0-100)")
	swapThreshold := flag.Float64("swap", 0, "Swap usage threshold percentage (0-100)")

	// Parse the command-line arguments
	flag.Parse()

	// Show version and exit if requested
	if *showVersion {
		fmt.Printf("Go System Monitor v%s\n", AppVersion)
		os.Exit(0)
	}

	// Show help and exit if requested
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Warning: Could not load configuration: %v. Using defaults.", err)
	}

	// Override configuration with command-line arguments if provided
	if *cpuThreshold > 0 {
		cfg.CPUThreshold = *cpuThreshold
	}
	if *memThreshold > 0 {
		cfg.MemoryThreshold = *memThreshold
	}
	if *diskThreshold > 0 {
		cfg.DiskThreshold = *diskThreshold
	}
	if *swapThreshold > 0 {
		cfg.SwapThreshold = *swapThreshold
	}

	fmt.Println("Go System Monitor Starting...")
	
	// Configure lipgloss for the terminal
	lipgloss.SetHasDarkBackground(true)
	
	// Run the Bubble Tea program
	p := tea.NewProgram(
		initialModel(cfg),
		tea.WithAltScreen(),       // Use the full terminal window
		tea.WithMouseCellMotion(), // Enable mouse support
	)
	
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
