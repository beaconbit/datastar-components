package piechart

import (
	"fmt"
	"math"
	"strings"
)

type Point struct {
	X float64
	Y float64
}

type Sector struct {
	Label      string
	Color      string
	Value      float64
	Percentage float64
}

type PieChartData struct {
	ID         string
	Title      string
	Sectors    []Sector
	Width      int
	Height     int
	RenderData []SectorRenderData // Precomputed for template
	SVG        string             // Full SVG string
	HTML       string             // Full chart HTML
}

func GetSectorPath(index int, sector Sector, sectors []Sector) string {
	// Calculate start and end angles based on cumulative percentages
	total := 0.0
	for i := 0; i < len(sectors); i++ {
		if i == index {
			break
		}
		total += sectors[i].Percentage
	}

	startAngle := total * 3.6 // Convert percentage to degrees (0-360)
	endAngle := startAngle + sector.Percentage*3.6

	// Convert to radians
	startRad := startAngle * math.Pi / 180
	endRad := endAngle * math.Pi / 180

	// Center and radius
	cx, cy, r := 200.0, 200.0, 150.0

	// Calculate points
	x1 := cx + r*math.Cos(startRad)
	y1 := cy + r*math.Sin(startRad)
	x2 := cx + r*math.Cos(endRad)
	y2 := cy + r*math.Sin(endRad)

	// Large arc flag: 1 if angle > 180 degrees
	largeArcFlag := 0
	if sector.Percentage > 50 {
		largeArcFlag = 1
	}

	// Create SVG path
	return fmt.Sprintf("M%.2f,%.2f L%.2f,%.2f A%.2f,%.2f 0 %d 1 %.2f,%.2f Z",
		cx, cy, x1, y1, r, r, largeArcFlag, x2, y2)
}

func GetLabelPosition(index int, sector Sector, sectors []Sector) Point {
	// Calculate middle angle of sector
	total := 0.0
	for i := 0; i < len(sectors); i++ {
		if i == index {
			break
		}
		total += sectors[i].Percentage
	}

	middleAngle := (total + sector.Percentage/2) * 3.6
	middleRad := middleAngle * math.Pi / 180

	// Position at 70% of radius
	cx, cy, r := 200.0, 200.0, 105.0

	return Point{
		X: cx + r*math.Cos(middleRad),
		Y: cy + r*math.Sin(middleRad),
	}
}

type SectorRenderData struct {
	Path           string
	LabelX         string
	LabelY         string
	Label          string
	Color          string
	Percentage     float64
	PercentageText string
}

func ComputeRenderData(data PieChartData) []SectorRenderData {
	result := make([]SectorRenderData, len(data.Sectors))
	for i, sector := range data.Sectors {
		pos := GetLabelPosition(i, sector, data.Sectors)
		result[i] = SectorRenderData{
			Path:           GetSectorPath(i, sector, data.Sectors),
			LabelX:         fmt.Sprintf("%.1f", pos.X),
			LabelY:         fmt.Sprintf("%.1f", pos.Y),
			Label:          sector.Label,
			Color:          sector.Color,
			Percentage:     sector.Percentage,
			PercentageText: FormatPercentage(sector.Percentage),
		}
	}
	return result
}

func GenerateSVGString(data PieChartData) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 400 400" xmlns="http://www.w3.org/2000/svg">`, data.Width, data.Height))
	sb.WriteString(`<circle cx="200" cy="200" r="150" fill="#f0f0f0" stroke="#ccc" stroke-width="1"/>`)

	renderData := ComputeRenderData(data)
	for _, rd := range renderData {
		sb.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" stroke="#fff" stroke-width="2"/>`, rd.Path, rd.Color))
		sb.WriteString(fmt.Sprintf(`<text x="%s" y="%s" text-anchor="middle" font-size="12" fill="#333">%s (%s)</text>`,
			rd.LabelX, rd.LabelY, rd.Label, rd.PercentageText))
	}

	sb.WriteString(`<circle cx="200" cy="200" r="50" fill="white"/>`)
	sb.WriteString(`<text x="200" y="200" text-anchor="middle" dy="5" font-size="14" font-weight="bold">Total</text>`)
	sb.WriteString(`</svg>`)
	return sb.String()
}

func GenerateChartHTML(data PieChartData) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<div id="%s" class="pie-chart-container">`, data.ID))
	sb.WriteString(fmt.Sprintf(`<h3>%s</h3>`, data.Title))
	sb.WriteString(`<div class="chart-svg">`)
	sb.WriteString(GenerateSVGString(data))
	sb.WriteString(`</div><div class="legend">`)

	for _, sector := range data.Sectors {
		sb.WriteString(`<div class="legend-item">`)
		sb.WriteString(fmt.Sprintf(`<span class="legend-color" style="background-color: %s"></span>`, sector.Color))
		sb.WriteString(fmt.Sprintf(`<span class="legend-label">%s</span>`, sector.Label))
		sb.WriteString(fmt.Sprintf(`<span class="legend-value">%s</span>`, FormatPercentage(sector.Percentage)))
		sb.WriteString(`</div>`)
	}
	sb.WriteString(`</div></div>`)
	return sb.String()
}

func FormatPercentage(p float64) string {
	return fmt.Sprintf("%.1f%%", p)
}
