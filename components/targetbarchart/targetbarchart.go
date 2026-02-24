package targetbarchart

import (
	"fmt"
	"strings"
)

// BarData represents a single bar in the target chart
type BarData struct {
	Label     string // Product label
	Target    int    // Total number of sections (target)
	Completed int    // Number of completed sections (green)
	Color     string // Base color for incomplete sections (optional)
}

// TargetBarChartData represents a target bar chart component's data
type TargetBarChartData struct {
	ID         string
	Title      string
	XAxisLabel string
	YAxisLabel string
	Bars       []BarData
	Width      int
	Height     int
	BarHeight  int    // Fixed height for each bar
	SVG        string // Full SVG string
	HTML       string // Full chart HTML
}

// DefaultTargetBarChart creates a target bar chart with default values
func DefaultTargetBarChart() TargetBarChartData {
	return TargetBarChartData{
		ID:         "target-bar-chart",
		Title:      "Target Progress Chart",
		XAxisLabel: "Products",
		YAxisLabel: "Sections",
		Width:      800,
		Height:     500,
		BarHeight:  200, // Fixed height for each bar
		Bars:       []BarData{},
	}
}

// GenerateSVGString generates the SVG for a target bar chart
func (b TargetBarChartData) GenerateSVGString() string {
	var sb strings.Builder

	// SVG dimensions and viewBox
	svgWidth := b.Width
	svgHeight := b.Height
	marginTop := 60
	marginRight := 20
	marginBottom := 80
	marginLeft := 80
	chartWidth := svgWidth - marginLeft - marginRight
	chartHeight := svgHeight - marginTop - marginBottom

	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		svgWidth, svgHeight, svgWidth, svgHeight))

	// Draw chart background
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="#f8f9fa" stroke="#ddd" stroke-width="1"/>`,
		marginLeft, marginTop, chartWidth, chartHeight))

	// Draw bars if we have data
	if len(b.Bars) > 0 {
		barAreaWidth := chartWidth / len(b.Bars)
		barWidth := barAreaWidth - 20 // Leave spacing between bars

		// Find max target for scaling Y axis (optional, since bar height fixed)
		maxTarget := 0
		for _, bar := range b.Bars {
			if bar.Target > maxTarget {
				maxTarget = bar.Target
			}
		}
		if maxTarget == 0 {
			maxTarget = 1
		}

		for i, bar := range b.Bars {
			// Calculate position for this bar
			x := marginLeft + (i * barAreaWidth) + (barAreaWidth-barWidth)/2
			y := marginTop + chartHeight - b.BarHeight

			// Draw bar background (empty bar)
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="#e9ecef" stroke="#ccc" stroke-width="1"/>`,
				x, y, barWidth, b.BarHeight))

			// Draw completed sections as green rectangles stacked from bottom
			if bar.Target > 0 {
				sectionHeight := float64(b.BarHeight) / float64(bar.Target)
				completed := bar.Completed
				if completed > bar.Target {
					completed = bar.Target
				}
				for s := 0; s < completed; s++ {
					sectionY := y + b.BarHeight - int(float64(s+1)*sectionHeight)
					sectionH := int(sectionHeight)
					// Ensure last section aligns with bottom
					if s == completed-1 {
						// Adjust height to avoid gaps due to rounding
						sectionH = int(float64(s+1)*sectionHeight) - int(float64(s)*sectionHeight)
					}
					sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="#10b981" stroke="#fff" stroke-width="0.5"/>`,
						x, sectionY, barWidth, sectionH))
				}
			}

			// Draw bar label (product name) under bar
			labelX := x + barWidth/2
			labelY := marginTop + chartHeight + 20
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="12" fill="#333">%s</text>`,
				labelX, labelY, bar.Label))

			// Draw progress label (completed/target) above bar
			progressText := fmt.Sprintf("%d/%d", bar.Completed, bar.Target)
			progressY := y - 10
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="11" fill="#333">%s</text>`,
				labelX, progressY, progressText))

			// Draw + and - buttons under each bar (as SVG groups for interactivity)
			buttonY := labelY + 25
			buttonSpacing := 25
			plusX := labelX - buttonSpacing
			minusX := labelX + buttonSpacing

			// Plus button
			sb.WriteString(fmt.Sprintf(`<g id="bar-%d-plus" class="target-bar-button" data-bar-index="%d" data-action="increment" data-on:click="@get('/api/targetbarchart/update?barIndex=%d&action=increment')" style="cursor: pointer;">`,
				i, i, i))
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="20" height="20" rx="3" fill="#4f46e5" stroke="#4f46e5" stroke-width="1"/>`,
				plusX-10, buttonY-10))
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="14" fill="white" dy="4">+</text>`,
				plusX, buttonY))
			sb.WriteString(`</g>`)

			// Minus button
			sb.WriteString(fmt.Sprintf(`<g id="bar-%d-minus" class="target-bar-button" data-bar-index="%d" data-action="decrement" data-on:click="@get('/api/targetbarchart/update?barIndex=%d&action=decrement')" style="cursor: pointer;">`,
				i, i, i))
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="20" height="20" rx="3" fill="#ef4444" stroke="#ef4444" stroke-width="1"/>`,
				minusX-10, buttonY-10))
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="14" fill="white" dy="4">-</text>`,
				minusX, buttonY))
			sb.WriteString(`</g>`)
		}

		// Draw Y axis with section ticks
		numYTicks := 5
		for i := 0; i <= numYTicks; i++ {
			value := (float64(i) / float64(numYTicks)) * float64(maxTarget)
			y := marginTop + chartHeight - int((float64(i)/float64(numYTicks))*float64(b.BarHeight))

			// Tick line
			sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#ccc" stroke-width="1"/>`,
				marginLeft-5, y, marginLeft, y))

			// Tick label
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="end" font-size="11" fill="#666">%.0f</text>`,
				marginLeft-10, y+4, value))
		}

		// Draw axes
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#333" stroke-width="2"/>`,
			marginLeft, marginTop, marginLeft, marginTop+chartHeight))
		sb.WriteString(fmt.Sprintf(`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="#333" stroke-width="2"/>`,
			marginLeft, marginTop+chartHeight, marginLeft+chartWidth, marginTop+chartHeight))

		// Chart title
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="16" font-weight="bold" fill="#333">%s</text>`,
			svgWidth/2, 30, b.Title))

		// X axis label
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="12" fill="#666">%s</text>`,
			svgWidth/2, svgHeight-20, b.XAxisLabel))

		// Y axis label
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="12" fill="#666" transform="rotate(-90,%d,%d)">%s</text>`,
			30, svgHeight/2, 30, svgHeight/2, b.YAxisLabel))
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}

// GenerateHTML generates the HTML for a target bar chart
func (b TargetBarChartData) GenerateHTML() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`<div id="%s" class="target-bar-chart-container">`, b.ID))
	sb.WriteString(fmt.Sprintf(`<h3>%s</h3>`, b.Title))
	sb.WriteString(`<div class="chart-svg">`)
	sb.WriteString(b.GenerateSVGString())
	sb.WriteString(`</div>`)

	// Legend (optional)
	if len(b.Bars) > 0 {
		sb.WriteString(`<div class="legend">`)
		sb.WriteString(`<div class="legend-item"><span class="legend-color" style="background-color: #10b981"></span><span class="legend-label">Completed sections</span></div>`)
		sb.WriteString(`<div class="legend-item"><span class="legend-color" style="background-color: #e9ecef"></span><span class="legend-label">Remaining sections</span></div>`)
		sb.WriteString(`</div>`)
	}

	sb.WriteString(`</div>`)
	return sb.String()
}

// WithBars creates a copy with specified bars
func (b TargetBarChartData) WithBars(bars []BarData) TargetBarChartData {
	b.Bars = bars
	return b
}

// WithTitle creates a copy with specified title
func (b TargetBarChartData) WithTitle(title string) TargetBarChartData {
	b.Title = title
	return b
}

// WithAxisLabels creates a copy with specified axis labels
func (b TargetBarChartData) WithAxisLabels(xLabel, yLabel string) TargetBarChartData {
	b.XAxisLabel = xLabel
	b.YAxisLabel = yLabel
	return b
}
