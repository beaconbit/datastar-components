package stackedbarchart

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/a-h/templ"
)

func renderComponentToString(c templ.Component) (string, error) {
	var buf strings.Builder
	ctx := context.Background()
	if err := c.Render(ctx, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Machine represents a washing machine with its delay data
type Machine struct {
	ID           int
	Name         string
	Color        string
	CurrentDelay int // Delay for current minute in seconds
	TotalDelay   int // Running total delay in seconds
}

// MinuteBar represents a single minute's stacked bar data
type MinuteBar struct {
	MinuteOffset  int // 0 = current minute, -1 = 1 min ago, ..., -9 = 9 min ago
	Timestamp     time.Time
	MachineDelays [3]int // Delays for each machine in seconds
	TotalDelay    int    // Sum of all machine delays for this minute
}

// StackedBarChartData represents the stacked bar chart component's data
type StackedBarChartData struct {
	ID          string
	Title       string
	XAxisLabel  string
	YAxisLabel  string
	Minutes     []MinuteBar // 10 minutes, index 0 is oldest (-9), index 9 is current (0)
	Machines    [3]Machine
	CurrentTime time.Time
	Width       int
	Height      int
	BarWidth    int // Width of each bar
	SVG         string
	HTML        string
	ClockHTML   string // Live clock display
	LegendHTML  string // Dynamic legend with buttons
}

// DefaultStackedBarChart creates a stacked bar chart with default values
func DefaultStackedBarChart() StackedBarChartData {
	now := time.Now()

	// Initialize machines
	machines := [3]Machine{
		{
			ID:           0,
			Name:         "Continuous Batch Washer 1",
			Color:        "#d8b4fe", // Light Purple
			CurrentDelay: 0,
			TotalDelay:   0,
		},
		{
			ID:           1,
			Name:         "Continuous Batch Washer 2",
			Color:        "#a855f7", // Medium Purple
			CurrentDelay: 0,
			TotalDelay:   0,
		},
		{
			ID:           2,
			Name:         "Continuous Batch Washer 3",
			Color:        "#7c3aed", // Dark Purple
			CurrentDelay: 0,
			TotalDelay:   0,
		},
	}

	// Initialize 10 minutes of data (9 minutes ago to current minute)
	minutes := make([]MinuteBar, 10)
	for i := 0; i < 10; i++ {
		offset := i - 9 // -9 (oldest) to 0 (current)
		timestamp := now.Add(time.Duration(offset) * time.Minute)
		minutes[i] = MinuteBar{
			MinuteOffset:  offset,
			Timestamp:     timestamp,
			MachineDelays: [3]int{0, 0, 0},
			TotalDelay:    0,
		}
	}

	return StackedBarChartData{
		ID:          "stacked-bar-chart",
		Title:       "Washing Machine Delay Monitor",
		XAxisLabel:  "Minutes Ago",
		YAxisLabel:  "Delay (seconds)",
		Minutes:     minutes,
		Machines:    machines,
		CurrentTime: now,
		Width:       1000,
		Height:      600,
		BarWidth:    80,
	}
}

// GenerateClockHTML generates the live clock display
func (s *StackedBarChartData) GenerateClockHTML() string {
	now := time.Now()
	currentTime := now.Format("15:04:05")
	component := Clock(currentTime, false)
	if html, err := renderComponentToString(component); err == nil {
		return html
	}
	return ""
}

// GenerateLegendHTML generates the dynamic legend with buttons and running counters
func (s *StackedBarChartData) GenerateLegendHTML() string {
	component := Legend(s.Machines)
	if html, err := renderComponentToString(component); err == nil {
		return html
	}
	return ""
}

// GenerateSVGString generates the SVG for the stacked bar chart
func (s *StackedBarChartData) GenerateSVGString() string {
	component := StackedBarChartSVG(*s)
	if html, err := renderComponentToString(component); err == nil {
		return html
	}
	return ""
}

// GenerateHTML generates the full HTML for the stacked bar chart component
func (s *StackedBarChartData) GenerateHTML() string {
	// Update component fields
	now := time.Now()
	currentTime := now.Format("15:04:05")
	clockComponent := Clock(currentTime, true)
	clockHTML, err := renderComponentToString(clockComponent)
	if err != nil {
		clockHTML = ""
	}
	s.ClockHTML = clockHTML
	s.LegendHTML = s.GenerateLegendHTML()
	s.SVG = s.GenerateSVGString()

	// Render full component using templ
	component := StackedBarChartComponent(*s, true)
	if html, err := renderComponentToString(component); err == nil {
		return html
	}
	return ""
}

// AddRandomDelay adds a random delay (1-15 seconds) to a machine's current minute
func (s *StackedBarChartData) AddRandomDelay(machineID int) {
	if machineID < 0 || machineID >= len(s.Machines) {
		return
	}

	// Add random delay between 1 and 15 seconds
	randomDelay := rand.Intn(15) + 1

	// Update current minute (last in array)
	currentMinute := &s.Minutes[len(s.Minutes)-1]
	currentMinute.MachineDelays[machineID] += randomDelay
	currentMinute.TotalDelay += randomDelay

	// Update machine totals
	s.Machines[machineID].CurrentDelay += randomDelay
	s.Machines[machineID].TotalDelay += randomDelay
}

// AdvanceMinute shifts the chart by one minute, dropping oldest, adding new current minute
func (s *StackedBarChartData) AdvanceMinute() {
	// Drop oldest minute (index 0) and shift remaining left
	// Create new slice with elements 1..9 (9 elements)
	// Update minute offsets for all bars (now indices 0..9 correspond to -9..0)
	for idx := range s.Minutes {
		s.Minutes[idx].MinuteOffset = idx - 1 // -9 to 0
	}
	newMinutes := s.Minutes[1:]

	// Create new current minute with zero delays
	now := time.Now()
	newCurrentMinute := MinuteBar{
		MinuteOffset:  0,
		Timestamp:     now,
		MachineDelays: [3]int{0, 0, 0},
		TotalDelay:    0,
	}

	// Append new current minute
	s.Minutes = append(newMinutes, newCurrentMinute)

	// Update current time
	s.CurrentTime = now

	// Reset machine CurrentDelay fields for new minute
	for i := range s.Machines {
		s.Machines[i].CurrentDelay = 0
	}

}
