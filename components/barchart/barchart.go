package barchart

import (
	"fmt"
	"strings"
)

// BarData represents a single bar in the chart
type BarData struct {
	Label string
	Value float64
	Color string
}

// BarChartData represents a bar chart component's data
type BarChartData struct {
	ID     string
	Title  string
	Bars   []BarData
	Width  int
	Height int
	SVG    string // Full SVG string
	HTML   string // Full chart HTML
}

// DefaultBarChart creates a bar chart with default values
func DefaultBarChart() BarChartData {
	return BarChartData{
		ID:     "bar-chart",
		Title:  "Bar Chart Demo",
		Width:  600,
		Height: 400,
		Bars:   []BarData{},
	}
}

// GenerateSVGString generates the SVG for a bar chart
func (b BarChartData) GenerateSVGString() string {
	var sb strings.Builder

	// SVG dimensions and viewBox
	svgWidth := b.Width
	svgHeight := b.Height
	marginTop := 40
	marginRight := 20
	marginBottom := 50
	marginLeft := 60
	chartWidth := svgWidth - marginLeft - marginRight
	chartHeight := svgHeight - marginTop - marginBottom

	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" xmlns="http://www.w3.org/2000/svg">`,
		svgWidth, svgHeight, svgWidth, svgHeight))

	// Find max value for scaling
	maxValue := 0.0
	for _, bar := range b.Bars {
		if bar.Value > maxValue {
			maxValue = bar.Value
		}
	}
	if maxValue == 0 {
		maxValue = 100 // Default max
	}

	// Draw chart background
	sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="#f8f9fa" stroke="#ddd" stroke-width="1"/>`,
		marginLeft, marginTop, chartWidth, chartHeight))

	// Draw bars if we have data
	if len(b.Bars) > 0 {
		barWidth := chartWidth/len(b.Bars) - 10
		barSpacing := 10

		for i, bar := range b.Bars {
			// Calculate bar height (percentage of chart height)
			barHeight := int((bar.Value / maxValue) * float64(chartHeight))

			// Bar position
			x := marginLeft + (i * (barWidth + barSpacing)) + barSpacing/2
			y := marginTop + chartHeight - barHeight

			// Draw bar
			sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" stroke="#fff" stroke-width="1">`,
				x, y, barWidth, barHeight, bar.Color))
			sb.WriteString(`<title>`)
			sb.WriteString(fmt.Sprintf("%s: %.1f", bar.Label, bar.Value))
			sb.WriteString(`</title>`)
			sb.WriteString(`</rect>`)

			// Draw bar label
			labelX := x + barWidth/2
			labelY := marginTop + chartHeight + 20
			sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="12" fill="#333">%s</text>`,
				labelX, labelY, bar.Label))

			// Draw value on top of bar
			if barHeight > 20 {
				valueY := y - 5
				sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="11" fill="#333">%.1f</text>`,
					labelX, valueY, bar.Value))
			}
		}

		// Draw Y axis labels
		numYTicks := 5
		for i := 0; i <= numYTicks; i++ {
			value := (float64(i) / float64(numYTicks)) * maxValue
			y := marginTop + chartHeight - int((float64(i)/float64(numYTicks))*float64(chartHeight))

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
			svgWidth/2, 25, b.Title))

		// Y axis label
		sb.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="middle" font-size="12" fill="#666" transform="rotate(-90,%d,%d)">Value</text>`,
			15, svgHeight/2, 15, svgHeight/2))
	}

	sb.WriteString(`</svg>`)
	return sb.String()
}

// GenerateHTML generates the HTML for a bar chart
func (b BarChartData) GenerateHTML() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`<div id="%s" class="bar-chart-container">`, b.ID))
	sb.WriteString(fmt.Sprintf(`<h3>%s</h3>`, b.Title))
	sb.WriteString(`<div class="chart-svg">`)
	sb.WriteString(b.GenerateSVGString())
	sb.WriteString(`</div>`)

	// Legend
	if len(b.Bars) > 0 {
		sb.WriteString(`<div class="legend">`)
		for _, bar := range b.Bars {
			sb.WriteString(`<div class="legend-item">`)
			sb.WriteString(fmt.Sprintf(`<span class="legend-color" style="background-color: %s"></span>`, bar.Color))
			sb.WriteString(fmt.Sprintf(`<span class="legend-label">%s</span>`, bar.Label))
			sb.WriteString(fmt.Sprintf(`<span class="legend-value">%.1f</span>`, bar.Value))
			sb.WriteString(`</div>`)
		}
		sb.WriteString(`</div>`)
	}

	sb.WriteString(`</div>`)
	return sb.String()
}

// WithBars creates a copy with specified bars
func (b BarChartData) WithBars(bars []BarData) BarChartData {
	b.Bars = bars
	return b
}

// WithTitle creates a copy with specified title
func (b BarChartData) WithTitle(title string) BarChartData {
	b.Title = title
	return b
}
