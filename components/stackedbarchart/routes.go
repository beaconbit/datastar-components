package stackedbarchart

import (
	"encoding/json"
	"log"
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
	log.Printf("GenerateInitialData: creating initial data")
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
	log.Printf("GenerateInitialData: data generated with %d minutes", len(data.Minutes))
	return data
}

// incrementHandler handles requests to add random delay to a machine
func (c *Component) incrementHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("incrementHandler: request received: %s %s", r.Method, r.URL.String())
	// Expect datastar request
	if r.URL.Query().Get("datastar") == "" {
		log.Printf("incrementHandler: missing datastar query param")
		http.Error(w, "datastar query param required", http.StatusBadRequest)
		return
	}

	// Read machine ID from query parameters
	machineIDStr := r.URL.Query().Get("machineId")
	if machineIDStr == "" {
		log.Printf("incrementHandler: missing machineId parameter")
		http.Error(w, "missing machineId parameter", http.StatusBadRequest)
		return
	}

	machineID, err := strconv.Atoi(machineIDStr)
	if err != nil || machineID < 0 || machineID >= 3 {
		log.Printf("incrementHandler: invalid machineId: %s", machineIDStr)
		http.Error(w, "invalid machineId", http.StatusBadRequest)
		return
	}

	log.Printf("incrementHandler: adding random delay to machine %d", machineID)
	// Add random delay to the specified machine
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data.AddRandomDelay(machineID)

	// Regenerate SVG and HTML
	log.Printf("incrementHandler: regenerating SVG and HTML")
	c.data.SVG = c.data.GenerateSVGString()
	c.data.HTML = c.data.GenerateHTML()

	// Create the stacked bar chart component with updated data
	chartComponent := StackedBarChart(c.data, true)

	// Create SSE response
	sse := datastar.NewSSE(w, r)

	// Render component to string
	var buf strings.Builder
	if err := chartComponent.Render(r.Context(), &buf); err != nil {
		log.Printf("incrementHandler: error rendering component: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := buf.String()
	log.Printf("incrementHandler: patching DOM with HTML length %d", len(html))
	// Patch the chart into the DOM
	sse.PatchElements(html)
}

// advanceMinuteHandler handles requests to advance the chart by one minute
func (c *Component) advanceMinuteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("advanceMinuteHandler: request received: %s %s", r.Method, r.URL.String())
	// Expect datastar request
	if r.URL.Query().Get("datastar") == "" {
		log.Printf("advanceMinuteHandler: missing datastar query param")
		http.Error(w, "datastar query param required", http.StatusBadRequest)
		return
	}

	log.Printf("advanceMinuteHandler: advancing chart minute")
	// Advance chart by one minute
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data.AdvanceMinute()

	// Regenerate SVG and HTML
	log.Printf("advanceMinuteHandler: regenerating SVG and HTML")
	c.data.SVG = c.data.GenerateSVGString()
	c.data.HTML = c.data.GenerateHTML()

	// Create the stacked bar chart component with updated data
	chartComponent := StackedBarChart(c.data, true)

	// Create SSE response
	sse := datastar.NewSSE(w, r)

	// Render component to string
	var buf strings.Builder
	if err := chartComponent.Render(r.Context(), &buf); err != nil {
		log.Printf("advanceMinuteHandler: error rendering component: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	html := buf.String()
	log.Printf("advanceMinuteHandler: patching DOM with HTML length %d", len(html))
	// Patch the chart into the DOM
	sse.PatchElements(html)
}

// clockTickHandler handles requests to update the live clock via SSE long-polling
func (c *Component) clockTickHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("clockTickHandler: request received: %s %s", r.Method, r.URL.String())

	// If not a datastar request, return plain HTML clock (backward compatibility)
	if r.URL.Query().Get("datastar") == "" {
		log.Printf("clockTickHandler: non-datastar request, returning plain HTML clock")
		clockHTML := c.data.GenerateClockHTML()
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(clockHTML))
		return
	}

	log.Printf("clockTickHandler: datastar request detected, setting up SSE")
	// Set up SSE connection
	sse := datastar.NewSSE(w, r)
	// Get flusher if available
	flusher, _ := w.(http.Flusher)

	// Panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("clockTickHandler panic recovered: %v", r)
		}
	}()

	// Get initial time
	now := time.Now()
	lastSecond := now.Second()
	lastMinute := now.Minute()
	log.Printf("clockTickHandler: initial time - minute:%d second:%d", lastMinute, lastSecond)

	// Track chart's current minute (for detecting when to advance)
	c.mu.RLock()
	chartMinute := c.data.CurrentTime.Minute()
	c.mu.RUnlock()
	log.Printf("clockTickHandler: chart minute: %d", chartMinute)

	// Create a ticker that fires every 100ms for time checking
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	log.Printf("clockTickHandler: entering SSE loop")
	for {
		select {
		case <-r.Context().Done():
			// Client disconnected
			log.Printf("clockTickHandler: client disconnected")
			return
		case <-ticker.C:
			now := time.Now()
			currentSecond := now.Second()
			currentMinute := now.Minute()

			// Check if wall-clock minute has changed
			if currentMinute != lastMinute {
				log.Printf("clockTickHandler: minute changed from %d to %d", lastMinute, currentMinute)
				lastMinute = currentMinute

				// Check if chart's minute is behind wall-clock minute
				c.mu.Lock()
				chartMinute = c.data.CurrentTime.Minute()
				log.Printf("clockTickHandler: BEFORE ADVANCE - chart minute: %d, wall-clock minute: %d, data.CurrentTime: %s",
					chartMinute, currentMinute, c.data.CurrentTime.Format("15:04:05"))

				// Log minutes array before advance
				log.Printf("clockTickHandler: BEFORE ADVANCE - Minutes array length: %d", len(c.data.Minutes))
				for i, minute := range c.data.Minutes {
					log.Printf("clockTickHandler: BEFORE minute[%d] - offset: %d, timestamp: %s, totalDelay: %d",
						i, minute.MinuteOffset, minute.Timestamp.Format("15:04:05"), minute.TotalDelay)
				}

				if chartMinute != currentMinute {
					log.Printf("clockTickHandler: chart minute (%d) behind wall-clock minute (%d), advancing", chartMinute, currentMinute)
					// BRANCH: Advance chart data and send entire modified page with new bar locations
					c.data.AdvanceMinute()

					// Log minutes array after advance
					log.Printf("clockTickHandler: AFTER ADVANCE - Minutes array length: %d, data.CurrentTime: %s",
						len(c.data.Minutes), c.data.CurrentTime.Format("15:04:05"))
					for i, minute := range c.data.Minutes {
						log.Printf("clockTickHandler: AFTER minute[%d] - offset: %d, timestamp: %s, totalDelay: %d",
							i, minute.MinuteOffset, minute.Timestamp.Format("15:04:05"), minute.TotalDelay)
					}

					c.data.SVG = c.data.GenerateSVGString()
					c.data.HTML = c.data.GenerateHTML()
					chartComponent := StackedBarChart(c.data, true)
					var buf strings.Builder
					if err := chartComponent.Render(r.Context(), &buf); err != nil {
						c.mu.Unlock()
						log.Printf("clockTickHandler: error rendering chart component: %v", err)
						// If rendering fails, abort the connection
						return
					}
					html := buf.String()
					c.mu.Unlock()
					sse.PatchElements(html)
					if flusher != nil {
						flusher.Flush()
					}
					log.Printf("clockTickHandler: sent updated chart HTML")
					// After advancing the chart, we've already sent a fresh clock;
					// no need to send a separate clock tick update this iteration.
					continue
				} else {
					log.Printf("clockTickHandler: chart minute (%d) already up to date with wall-clock minute (%d), no advance needed", chartMinute, currentMinute)
				}
				c.mu.Unlock()
			}

			// Check if wall-clock second has changed
			if currentSecond != lastSecond {
				log.Printf("clockTickHandler: second changed from %d to %d", lastSecond, currentSecond)
				lastSecond = currentSecond
				// BRANCH: Simple clock tick update – send only the new clock HTML
				c.mu.RLock()
				clockHTML := c.data.GenerateClockHTML()
				c.mu.RUnlock()
				log.Printf("clockTickHandler: clock HTML length %d", len(clockHTML))
				sse.PatchElements(clockHTML)
				if flusher != nil {
					flusher.Flush()
				}
				log.Printf("clockTickHandler: sent clock tick update")
			}
		}
	}
}

// randomizeHandler handles requests to randomize all data
func (c *Component) randomizeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("randomizeHandler: request received: %s %s", r.Method, r.URL.String())
	// Check if this is a Datastar request
	if r.URL.Query().Get("datastar") != "" {
		log.Printf("randomizeHandler: datastar request detected")
		// Read signals from request (ignore errors)
		var signals map[string]interface{}
		datastar.ReadSignals(r, &signals) // Ignore error for GET requests

		// Generate new random data
		c.mu.Lock()
		defer c.mu.Unlock()
		log.Printf("randomizeHandler: generating new random data")
		c.data = c.GenerateInitialData()

		// Create the stacked bar chart component
		chartComponent := StackedBarChart(c.data, true)

		// Create SSE response
		sse := datastar.NewSSE(w, r)

		// Render component to string
		var buf strings.Builder
		if err := chartComponent.Render(r.Context(), &buf); err != nil {
			log.Printf("randomizeHandler: error rendering component: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		html := buf.String()
		log.Printf("randomizeHandler: patching DOM with HTML length %d", len(html))
		// Patch the chart into the DOM
		sse.PatchElements(html)

		return
	}

	log.Printf("randomizeHandler: non-datastar request, returning plain HTML")
	// Fallback: return just the chart HTML for non-Datastar requests
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = c.GenerateInitialData()
	chartComponent := StackedBarChart(c.data, true)

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
