package stackedbarchart

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/starfederation/datastar-go/datastar"
)

// Component implements the stacked bar chart component
type Component struct {
	rand *rand.Rand
	mu   sync.RWMutex
	data StackedBarChartData
}

// New creates a new stacked bar chart component instance
func New() *Component {
	now := time.Now()
	comp := &Component{
		rand: rand.New(rand.NewSource(now.UnixNano())),
	}
	comp.data = comp.GenerateInitialData()
	return comp
}

// GenerateInitialData creates initial stacked bar chart data
func (c *Component) GenerateInitialData() StackedBarChartData {
	data := DefaultStackedBarChart()

	// Add some random delays for demonstration
	// Add random delays to current minute for each machine
	for i := 0; i < 3; i++ {
		delay := c.rand.Intn(30) + 1
		data.AddRandomDelay(i)
		// Set a specific delay for demo
		data.Machines[i].TotalDelay = delay * 2
	}

	// Add some historical delays
	for i := 0; i < len(data.Minutes)-1; i++ {
		for j := 0; j < 3; j++ {
			delay := c.rand.Intn(20)
			data.Minutes[i].MachineDelays[j] = delay
			data.Minutes[i].TotalDelay += delay
		}
	}

	data.SVG = data.GenerateSVGString()
	data.HTML = data.GenerateHTML()
	return data
}

// incrementHandler handles requests to add random delay to a machine
func (c *Component) incrementHandler(w http.ResponseWriter, r *http.Request) {
	// Expect datastar request
	if r.URL.Query().Get("datastar") == "" {
		http.Error(w, "datastar query param required", http.StatusBadRequest)
		return
	}

	// Read machine ID from query parameters
	machineIDStr := r.URL.Query().Get("machineId")
	if machineIDStr == "" {
		http.Error(w, "missing machineId parameter", http.StatusBadRequest)
		return
	}

	machineID, err := strconv.Atoi(machineIDStr)
	if err != nil || machineID < 0 || machineID >= 3 {
		http.Error(w, "invalid machineId", http.StatusBadRequest)
		return
	}

	// Add random delay to the specified machine
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data.AddRandomDelay(machineID)

	// Regenerate SVG and HTML
	c.data.SVG = c.data.GenerateSVGString()
	c.data.HTML = c.data.GenerateHTML()

	// Create the stacked bar chart component with updated data
	chartComponent := StackedBarChart(c.data, false)

	// Create SSE response
	sse := datastar.NewSSE(w, r)

	// Render component to string
	var buf strings.Builder
	if err := chartComponent.Render(r.Context(), &buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Patch the chart into the DOM
	sse.PatchElements(buf.String())
}

// advanceMinuteHandler handles requests to advance the chart by one minute
func (c *Component) advanceMinuteHandler(w http.ResponseWriter, r *http.Request) {
	// Expect datastar request
	if r.URL.Query().Get("datastar") == "" {
		http.Error(w, "datastar query param required", http.StatusBadRequest)
		return
	}

	// Advance chart by one minute
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data.AdvanceMinute()

	// Regenerate SVG and HTML
	c.data.SVG = c.data.GenerateSVGString()
	c.data.HTML = c.data.GenerateHTML()

	// Create the stacked bar chart component with updated data
	chartComponent := StackedBarChart(c.data, false)

	// Create SSE response
	sse := datastar.NewSSE(w, r)

	// Render component to string
	var buf strings.Builder
	if err := chartComponent.Render(r.Context(), &buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Patch the chart into the DOM
	sse.PatchElements(buf.String())
}

// clockTickHandler handles requests to update the live clock via SSE long-polling
func (c *Component) clockTickHandler(w http.ResponseWriter, r *http.Request) {
	// If not a datastar request, return plain HTML clock (backward compatibility)
	if r.URL.Query().Get("datastar") == "" {
		clockHTML := c.data.GenerateClockHTML()
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(clockHTML))
		return
	}

	// Set up SSE connection
	sse := datastar.NewSSE(w, r)

	// Get initial time
	now := time.Now()
	lastSecond := now.Second()
	lastMinute := now.Minute()

	// Track chart's current minute (for detecting when to advance)
	c.mu.RLock()
	chartMinute := c.data.CurrentTime.Minute()
	c.mu.RUnlock()

	// Create a ticker that fires every 100ms
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			// Client disconnected
			return
		case <-ticker.C:
			now := time.Now()
			currentSecond := now.Second()
			currentMinute := now.Minute()

			// Check if wall-clock minute has changed
			if currentMinute != lastMinute {
				lastMinute = currentMinute

				// Check if chart's minute is behind wall-clock minute
				c.mu.Lock()
				chartMinute = c.data.CurrentTime.Minute()
				if chartMinute != currentMinute {
					// BRANCH: Advance chart data and send entire modified page with new bar locations
					c.data.AdvanceMinute()
					c.data.SVG = c.data.GenerateSVGString()
					c.data.HTML = c.data.GenerateHTML()
					chartComponent := StackedBarChart(c.data, false)
					var buf strings.Builder
					if err := chartComponent.Render(r.Context(), &buf); err != nil {
						c.mu.Unlock()
						// If rendering fails, abort the connection
						return
					}
					html := buf.String()
					c.mu.Unlock()
					sse.PatchElements(html)
					// After advancing the chart, we've already sent a fresh clock;
					// no need to send a separate clock tick update this iteration.
					continue
				}
				c.mu.Unlock()
			}

			// Check if wall-clock second has changed
			if currentSecond != lastSecond {
				lastSecond = currentSecond
				// BRANCH: Simple clock tick update – send only the new clock HTML
				clockHTML := c.data.GenerateClockHTML()
				sse.PatchElements(clockHTML)
			}
		}
	}
}

// randomizeHandler handles requests to randomize all data
func (c *Component) randomizeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if this is a Datastar request
	if r.URL.Query().Get("datastar") != "" {
		// Read signals from request (ignore errors)
		var signals map[string]interface{}
		datastar.ReadSignals(r, &signals) // Ignore error for GET requests

		// Generate new random data
		c.mu.Lock()
		defer c.mu.Unlock()
		c.data = c.GenerateInitialData()

		// Create the stacked bar chart component
		chartComponent := StackedBarChart(c.data, false)

		// Create SSE response
		sse := datastar.NewSSE(w, r)

		// Render component to string
		var buf strings.Builder
		if err := chartComponent.Render(r.Context(), &buf); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Patch the chart into the DOM
		sse.PatchElements(buf.String())

		return
	}

	// Fallback: return just the chart HTML for non-Datastar requests
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = c.GenerateInitialData()
	chartComponent := StackedBarChart(c.data, false)

	w.Header().Set("Content-Type", "text/html")
	chartComponent.Render(r.Context(), w)
}

// signalsHandler handles requests for stacked bar chart signals
func (c *Component) signalsHandler(w http.ResponseWriter, r *http.Request) {
	// Example endpoint to patch signals
	sse := datastar.NewSSE(w, r)

	c.mu.RLock()
	defer c.mu.RUnlock()
	// Send current stats
	signals := map[string]interface{}{
		"lastUpdated":  time.Now().Format(time.RFC3339),
		"totalMinutes": len(c.data.Minutes),
		"currentTime":  c.data.CurrentTime.Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(signals)
	sse.PatchSignals(jsonData)
}

// RegisterRoutes registers HTTP routes for the stacked bar chart component
func (c *Component) RegisterRoutes(r chi.Router) {
	r.Get("/increment", c.incrementHandler)
	r.Get("/advance", c.advanceMinuteHandler)
	r.Get("/tick", c.clockTickHandler)
	r.Get("/randomize", c.randomizeHandler)
	r.Get("/signals", c.signalsHandler)
}

// RegisterStatic registers static asset routes for the stacked bar chart component
func (c *Component) RegisterStatic(r chi.Router) {
	// Stacked bar chart doesn't have static assets yet
	// r.Handle("/assets/*", http.FileServer(http.Dir("./components/stackedbarchart/assets")))
}
