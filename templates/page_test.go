package templates

import (
	"bytes"
	"context"
	"piechart-demo/components"
	"strings"
	"testing"
)

func TestPageButton(t *testing.T) {
	colors := []string{"#4f46e5", "#10b981"}
	labels := []string{"Tech", "Health"}
	sectors := []components.Sector{
		{Label: labels[0], Color: colors[0], Value: 50, Percentage: 50},
		{Label: labels[1], Color: colors[1], Value: 50, Percentage: 50},
	}
	data := components.PieChartData{
		ID:      "pie-chart",
		Title:   "Test",
		Sectors: sectors,
		Width:   500,
		Height:  500,
	}
	data.RenderData = components.ComputeRenderData(data)
	data.SVG = components.GenerateSVGString(data)
	data.HTML = components.GenerateChartHTML(data)

	var buf bytes.Buffer
	component := Page(data)
	component.Render(context.Background(), &buf)

	html := buf.String()
	if !strings.Contains(html, `data-on:click="@get('/test/randomize')"`) {
		t.Errorf("Button attribute incorrect. HTML snippet:\n%s", html)
	}
}
