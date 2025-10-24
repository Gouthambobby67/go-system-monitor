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

		// Card visibility toggles (only in All view)
		case "s", "c", "m", "d", "n", "a", "g":
			if m.dashboard.ActiveTab() == 7 { // All/Combined view
				m.dashboard.UpdateCardConfig(msg.String())
				return m, nil
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