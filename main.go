package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"piechart-demo/components"
	"piechart-demo/templates"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/starfederation/datastar-go/datastar"
)

var colors = []string{
	"#4f46e5", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6",
	"#06b6d4", "#84cc16", "#f97316", "#ec4899", "#14b8a6",
}

var labels = []string{
	"Technology", "Healthcare", "Finance", "Energy", "Consumer",
	"Industrials", "Materials", "Real Estate", "Utilities", "Communications",
}

type Server struct {
	rand *rand.Rand
}

func NewServer() *Server {
	return &Server{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Server) generateRandomData() components.PieChartData {
	// Generate 5-8 random sectors
	numSectors := s.rand.Intn(4) + 5

	sectors := make([]components.Sector, numSectors)
	totalValue := 0.0

	// Generate random values
	for i := 0; i < numSectors; i++ {
		value := float64(s.rand.Intn(100) + 10)
		sectors[i] = components.Sector{
			Label: labels[i%len(labels)],
			Color: colors[i%len(colors)],
			Value: value,
		}
		totalValue += value
	}

	// Calculate percentages
	for i := range sectors {
		sectors[i].Percentage = (sectors[i].Value / totalValue) * 100
	}

	data := components.PieChartData{
		ID:      "pie-chart",
		Title:   "Market Distribution",
		Sectors: sectors,
		Width:   500,
		Height:  500,
	}

	// Compute render data, SVG and full HTML
	data.RenderData = components.ComputeRenderData(data)
	data.SVG = components.GenerateSVGString(data)
	data.HTML = components.GenerateChartHTML(data)
	return data
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	data := s.generateRandomData()
	component := templates.Page(data)

	w.Header().Set("Content-Type", "text/html")

	// Debug: log button attribute
	if r.URL.Query().Get("debug") == "1" {
		var buf strings.Builder
		component.Render(r.Context(), &buf)
		html := buf.String()
		idx := strings.Index(html, "data-on:click=")
		if idx >= 0 {
			end := strings.Index(html[idx:], ">")
			if end >= 0 {
				log.Printf("DEBUG button attribute: %s", html[idx:idx+end])
			}
		}
		w.Write([]byte(html))
	} else {
		component.Render(r.Context(), w)
	}
}

func (s *Server) randomizeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if this is a Datastar request by looking for the datastar query param
	if r.URL.Query().Get("datastar") != "" {
		// Read signals from request (though we don't need them for this demo)
		var signals map[string]interface{}
		if err := datastar.ReadSignals(r, &signals); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Generate new random data
		data := s.generateRandomData()

		// Create the pie chart component
		chartComponent := components.PieChart(data)

		// Create SSE response
		sse := datastar.NewSSE(w, r)

		// Render component to string
		var buf strings.Builder
		if err := chartComponent.Render(r.Context(), &buf); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Patch the chart into the DOM (morphing into element with id="pie-chart")
		sse.PatchElements(buf.String())

		return
	}

	// Fallback: return just the chart HTML for non-Datastar requests
	data := s.generateRandomData()
	chartComponent := components.PieChart(data)

	w.Header().Set("Content-Type", "text/html")
	chartComponent.Render(r.Context(), w)
}

func (s *Server) signalsHandler(w http.ResponseWriter, r *http.Request) {
	// Example endpoint to patch signals
	sse := datastar.NewSSE(w, r)

	// Send some example signals
	signals := map[string]interface{}{
		"lastUpdated":  time.Now().Format(time.RFC3339),
		"totalSectors": 7,
	}

	jsonData, _ := json.Marshal(signals)
	sse.PatchSignals(jsonData)
}

// Middleware to log request duration
func requestTimer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s completed in %v", r.Method, r.URL.Path, duration)
	})
}

func main() {
	server := NewServer()

	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(requestTimer)

	// Routes
	r.Get("/", server.homeHandler)
	r.Get("/randomize", server.randomizeHandler)
	r.Get("/signals", server.signalsHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("Press Ctrl+C to stop")

	log.Fatal(http.ListenAndServe(port, r))
}
