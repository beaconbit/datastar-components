package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"piechart-demo/components/barchart"
	"piechart-demo/components/piechart"
	"piechart-demo/components/stackedbarchart"
	"piechart-demo/components/targetbarchart"
	"piechart-demo/templates"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// homeHandler serves the component library home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	component := templates.Home()
	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
}

// piechartPageHandler serves a page with just the pie chart component
func piechartPageHandler(w http.ResponseWriter, r *http.Request) {
	piechartComp := piechart.New()
	data := piechartComp.GenerateRandomData()

	// Create a simple page with just the pie chart
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Pie Chart Component</title>
		<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js"></script>
		<style>
			body {
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
				margin: 0;
				padding: 40px;
				background-color: #f5f5f5;
				display: flex;
				flex-direction: column;
				align-items: center;
				min-height: 100vh;
			}
			.container {
				background: white;
				border-radius: 12px;
				padding: 40px;
				box-shadow: 0 4px 12px rgba(0,0,0,0.1);
				max-width: 800px;
				width: 100%;
			}
			h1 {
				color: #333;
				margin-bottom: 8px;
			}
			.subtitle {
				color: #666;
				margin-bottom: 24px;
			}
			.controls {
				margin-bottom: 24px;
				display: flex;
				gap: 12px;
			}
			button {
				background: #4f46e5;
				color: white;
				border: none;
				padding: 12px 24px;
				border-radius: 8px;
				font-size: 16px;
				font-weight: 600;
				cursor: pointer;
				transition: background 0.2s;
			}
			button:hover {
				background: #4338ca;
			}
			.back-link {
				margin-top: 24px;
				color: #4f46e5;
				text-decoration: none;
				font-weight: 500;
			}
			.back-link:hover {
				text-decoration: underline;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Pie Chart Component</h1>
			<p class="subtitle">Standalone component demonstration</p>
			
			<div class="controls">
				<button data-on:click="@get('/api/piechart/randomize')">
					Randomize Chart
				</button>
			</div>
			
			<div class="chart-container">
				` + data.HTML + `
			</div>
			
			<a href="/" class="back-link">← Back to Component Library</a>
		</div>
	</body>
	</html>`))
}

// barchartPageHandler serves a page with just the bar chart component
func barchartPageHandler(w http.ResponseWriter, r *http.Request) {
	barchartComp := barchart.New()
	data := barchartComp.GenerateRandomData()

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Bar Chart Component</title>
		<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js"></script>
		<style>
			body {
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
				margin: 0;
				padding: 40px;
				background-color: #f5f5f5;
				display: flex;
				flex-direction: column;
				align-items: center;
				min-height: 100vh;
			}
			.container {
				background: white;
				border-radius: 12px;
				padding: 40px;
				box-shadow: 0 4px 12px rgba(0,0,0,0.1);
				max-width: 800px;
				width: 100%;
			}
			h1 {
				color: #333;
				margin-bottom: 8px;
			}
			.subtitle {
				color: #666;
				margin-bottom: 24px;
			}
			.controls {
				margin-bottom: 24px;
				display: flex;
				gap: 12px;
			}
			button {
				background: #10b981;
				color: white;
				border: none;
				padding: 12px 24px;
				border-radius: 8px;
				font-size: 16px;
				font-weight: 600;
				cursor: pointer;
				transition: background 0.2s;
			}
			button:hover {
				background: #0da271;
			}
			.back-link {
				margin-top: 24px;
				color: #4f46e5;
				text-decoration: none;
				font-weight: 500;
			}
			.back-link:hover {
				text-decoration: underline;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Bar Chart Component</h1>
			<p class="subtitle">Standalone component demonstration</p>
			
			<div class="controls">
				<button data-on:click="@get('/api/barchart/randomize')">
					Randomize Chart
				</button>
			</div>
			
			<div class="chart-container">
				` + data.HTML + `
			</div>
			
			<a href="/" class="back-link">← Back to Component Library</a>
		</div>
	</body>
	</html>`))
}

// buttonPageHandler serves a page with just the button component
func buttonPageHandler(w http.ResponseWriter, r *http.Request) {
	// buttonComp := button.New() // Not used in static HTML

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Button Component</title>
		<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js"></script>
		<style>
			body {
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
				margin: 0;
				padding: 40px;
				background-color: #f5f5f5;
				display: flex;
				flex-direction: column;
				align-items: center;
				min-height: 100vh;
			}
			.container {
				background: white;
				border-radius: 12px;
				padding: 40px;
				box-shadow: 0 4px 12px rgba(0,0,0,0.1);
				max-width: 800px;
				width: 100%;
			}
			h1 {
				color: #333;
				margin-bottom: 8px;
			}
			.subtitle {
				color: #666;
				margin-bottom: 24px;
			}
			.controls {
				margin-bottom: 24px;
				display: flex;
				flex-wrap: wrap;
				gap: 12px;
			}
			button {
				background: #f59e0b;
				color: white;
				border: none;
				padding: 12px 24px;
				border-radius: 8px;
				font-size: 16px;
				font-weight: 600;
				cursor: pointer;
				transition: background 0.2s;
			}
			button:hover {
				background: #d4890a;
			}
			.back-link {
				margin-top: 24px;
				color: #4f46e5;
				text-decoration: none;
				font-weight: 500;
			}
			.back-link:hover {
				text-decoration: underline;
			}
			.component-demo {
				background: #f8f9fa;
				border-radius: 8px;
				padding: 24px;
				margin-top: 24px;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Button Component</h1>
			<p class="subtitle">Standalone component demonstration</p>
			
			<div class="controls">
				<button data-on:click="@get('/api/button/click')">
					Click Me
				</button>
				<button data-on:click="@get('/api/button/toggle')">
					Toggle State
				</button>
			</div>
			
			<div class="component-demo">
				<div id="demo-button">
					<button id="demo-button" class="btn btn-primary btn-medium" data-click-count="0">
						Click me
					</button>
				</div>
			</div>
			
			<a href="/" class="back-link">← Back to Component Library</a>
		</div>
	</body>
	</html>`))
}

// formPageHandler serves a page with just the form component
func formPageHandler(w http.ResponseWriter, r *http.Request) {
	// formComp := form.New() // Not used in static HTML

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Form Component</title>
		<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js"></script>
		<style>
			body {
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
				margin: 0;
				padding: 40px;
				background-color: #f5f5f5;
				display: flex;
				flex-direction: column;
				align-items: center;
				min-height: 100vh;
			}
			.container {
				background: white;
				border-radius: 12px;
				padding: 40px;
				box-shadow: 0 4px 12px rgba(0,0,0,0.1);
				max-width: 800px;
				width: 100%;
			}
			h1 {
				color: #333;
				margin-bottom: 8px;
			}
			.subtitle {
				color: #666;
				margin-bottom: 24px;
			}
			.controls {
				margin-bottom: 24px;
				display: flex;
				gap: 12px;
			}
			button {
				background: #ef4444;
				color: white;
				border: none;
				padding: 12px 24px;
				border-radius: 8px;
				font-size: 16px;
				font-weight: 600;
				cursor: pointer;
				transition: background 0.2s;
			}
			button:hover {
				background: #dc2626;
			}
			.back-link {
				margin-top: 24px;
				color: #4f46e5;
				text-decoration: none;
				font-weight: 500;
			}
			.back-link:hover {
				text-decoration: underline;
			}
			.component-demo {
				background: #f8f9fa;
				border-radius: 8px;
				padding: 24px;
				margin-top: 24px;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Form Component</h1>
			<p class="subtitle">Standalone component demonstration</p>
			
			<div class="controls">
				<button data-on:click="@get('/api/form/submit')">
					Submit Form
				</button>
				<button data-on:click="@get('/api/form/reset')">
					Reset Form
				</button>
			</div>
			
			<div class="component-demo">
				<div id="demo-form">
					<form id="demo-form" class="form" action="/api/form/submit" method="POST">
						<h3>Contact Form</h3>
						<div class="form-field">
							<label for="name">Name</label>
							<input type="text" id="name" name="name" value="" placeholder="Enter your name" required>
						</div>
						<div class="form-field">
							<label for="email">Email</label>
							<input type="email" id="email" name="email" value="" placeholder="Enter your email" required>
						</div>
						<div class="form-field">
							<label for="message">Message</label>
							<textarea id="message" name="message" placeholder="Enter your message" required></textarea>
						</div>
						<div class="form-field">
							<label for="newsletter">Subscribe to newsletter</label>
							<input type="checkbox" id="newsletter" name="newsletter" value="true" checked>
						</div>
						<div class="form-actions">
							<button type="submit" class="btn btn-primary">Submit</button>
							<button type="reset" class="btn btn-secondary">Reset</button>
						</div>
					</form>
				</div>
			</div>
			
			<a href="/" class="back-link">← Back to Component Library</a>
		</div>
	</body>
	</html>`))
}

// targetbarchartPageHandler serves a page with just the target bar chart component
func targetbarchartPageHandler(w http.ResponseWriter, r *http.Request) {
	targetbarchartComp := targetbarchart.New()
	data := targetbarchartComp.GenerateRandomData()

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Target Bar Chart Component</title>
		<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js"></script>
		<style>
			body {
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
				margin: 0;
				padding: 40px;
				background-color: #f5f5f5;
				display: flex;
				flex-direction: column;
				align-items: center;
				min-height: 100vh;
			}
			.container {
				background: white;
				border-radius: 12px;
				padding: 40px;
				box-shadow: 0 4px 12px rgba(0,0,0,0.1);
				max-width: 1000px;
				width: 100%;
			}
			h1 {
				color: #333;
				margin-bottom: 8px;
			}
			.subtitle {
				color: #666;
				margin-bottom: 24px;
			}
			.controls {
				margin-bottom: 24px;
				display: flex;
				gap: 12px;
			}
			button {
				background: #8b5cf6;
				color: white;
				border: none;
				padding: 12px 24px;
				border-radius: 8px;
				font-size: 16px;
				font-weight: 600;
				cursor: pointer;
				transition: background 0.2s;
			}
			button:hover {
				background: #7c3aed;
			}
			.back-link {
				margin-top: 24px;
				color: #4f46e5;
				text-decoration: none;
				font-weight: 500;
			}
			.back-link:hover {
				text-decoration: underline;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Target Bar Chart Component</h1>
			<p class="subtitle">Standalone component demonstration</p>
			
			<div class="controls">
				<button data-on:click="@get('/api/targetbarchart/randomize')">
					Randomize Chart
				</button>
			</div>
			
			<div class="chart-container">
				` + data.HTML + `
			</div>
			
			<a href="/" class="back-link">← Back to Component Library</a>
		</div>
	</body>
	</html>`))
}

// stackedbarchartPageHandler serves a page with just the stacked bar chart component
func stackedbarchartPageHandler(w http.ResponseWriter, r *http.Request) {
	stackedbarchartComp := stackedbarchart.New()
	data := stackedbarchartComp.GenerateInitialData()

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Stacked Bar Chart Component</title>
		<script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js"></script>
		<style>
			body {
				font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
				margin: 0;
				padding: 40px;
				background-color: #f5f5f5;
				display: flex;
				flex-direction: column;
				align-items: center;
				min-height: 100vh;
			}
			.container {
				background: white;
				border-radius: 12px;
				padding: 40px;
				box-shadow: 0 4px 12px rgba(0,0,0,0.1);
				max-width: 1200px;
				width: 100%;
			}
			h1 {
				color: #333;
				margin-bottom: 8px;
			}
			.subtitle {
				color: #666;
				margin-bottom: 24px;
			}
			.controls {
				margin-bottom: 24px;
				display: flex;
				gap: 12px;
				flex-wrap: wrap;
			}
			button {
				background: #8b5cf6;
				color: white;
				border: none;
				padding: 12px 24px;
				border-radius: 8px;
				font-size: 16px;
				font-weight: 600;
				cursor: pointer;
				transition: background 0.2s;
			}
			button:hover {
				background: #7c3aed;
			}
			.back-link {
				margin-top: 24px;
				color: #a855f7;
				text-decoration: none;
				font-weight: 500;
			}
			.back-link:hover {
				text-decoration: underline;
			}
			.live-clock {
				background: #f8f9fa;
				border-radius: 8px;
				padding: 16px;
				margin-bottom: 20px;
				text-align: center;
				border: 2px solid #e9ecef;
			}
			.clock-label {
				font-size: 14px;
				color: #666;
				margin-bottom: 8px;
			}
			.clock-time {
				font-size: 24px;
				font-weight: bold;
				color: #333;
				font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
			}
			.stacked-legend {
				background: #f8f9fa;
				border-radius: 8px;
				padding: 20px;
				margin-bottom: 24px;
				border: 1px solid #e9ecef;
			}
			.legend-title {
				font-size: 16px;
				font-weight: 600;
				color: #333;
				margin-bottom: 12px;
			}
			.legend-items {
				display: flex;
				flex-direction: column;
				gap: 10px;
			}
			.legend-item {
				display: flex;
				align-items: center;
				gap: 12px;
			}
			.legend-color {
				width: 20px;
				height: 20px;
				border-radius: 4px;
				flex-shrink: 0;
			}
			.legend-label {
				flex-grow: 1;
				font-size: 14px;
				color: #333;
			}
			.legend-delay {
				font-size: 14px;
				font-weight: 600;
				color: #333;
				font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
				min-width: 120px;
			}
			.legend-button {
				background: #a855f7;
				color: white;
				border: none;
				padding: 6px 12px;
				border-radius: 4px;
				font-size: 14px;
				cursor: pointer;
				transition: background 0.2s;
			}
			.legend-button:hover {
				background: #9333ea;
			}
			.chart-svg {
				margin-top: 20px;
				border-radius: 8px;
				overflow: hidden;
				border: 1px solid #e9ecef;
			}
			@keyframes pulse {
				0% { opacity: 1; }
				50% { opacity: 0.7; }
				100% { opacity: 1; }
			}
			.current-minute {
				animation: pulse 2s infinite;
			}
		</style>
		<script>
			// Clock update system with server fallback
			(function() {
				let lastServerSuccess = Date.now();
				const CLOCK_TIMEOUT = 3000; // 3 seconds
				const CLOCK_ELEMENT_ID = 'stacked-chart-clock';
				let clockUpdateInterval = null;
				let serverUpdateInterval = null;
				
				// Function to update clock display with HH:MM:SS format
				function updateClockDisplay(timeString) {
					const clockElement = document.getElementById(CLOCK_ELEMENT_ID);
					if (!clockElement) return;
					
					// Find the clock-time element within the clock container
					const timeElement = clockElement.querySelector('.clock-time');
					if (timeElement) {
						timeElement.textContent = timeString;
					}
				}
				
				// Client-side clock update using local time
				function updateClockClient() {
					const now = new Date();
					const timeString = now.toTimeString().split(' ')[0]; // HH:MM:SS
					updateClockDisplay(timeString);
				}
				
				// Server clock update attempt - only send datastar=true at second 00
				function attemptServerClockUpdate() {
					const now = new Date();
					const seconds = now.getSeconds();
					// Only send datastar=true when seconds are 00 (start of minute)
					// For other seconds, send datastar=false to reduce SSE overhead
					const datastarParam = seconds === 0 ? 'true' : 'false';
					
					fetch('/api/stackedbarchart/tick?datastar=' + datastarParam)
						.then(response => {
							if (response.ok) {
								lastServerSuccess = Date.now();
								return;
							}
							console.error('Clock tick failed');
						})
						.catch(err => {
							console.error('Clock tick error:', err);
						});
				}
				
				// Monitor clock update status and attempt server updates
				serverUpdateInterval = setInterval(() => {
					const now = Date.now();
					const timeSinceLastSuccess = now - lastServerSuccess;
					
					// Always attempt server update (regardless of mode)
					// This allows us to detect when server comes back online
					attemptServerClockUpdate();
					
					// Log warning if server hasn't responded in 3 seconds
					if (timeSinceLastSuccess > CLOCK_TIMEOUT) {
						console.warn('Server not responding for', Math.round(timeSinceLastSuccess / 1000), 'seconds');
					}
				}, 1000);
				
				// Start client-side clock updates immediately (runs every second)
				updateClockClient();
				clockUpdateInterval = setInterval(updateClockClient, 1000);
				
				// Cleanup on page unload
				window.addEventListener('beforeunload', () => {
					if (clockUpdateInterval) clearInterval(clockUpdateInterval);
					if (serverUpdateInterval) clearInterval(serverUpdateInterval);
				});
				
				// Initial server update
				attemptServerClockUpdate();
			})();

			// Auto-advance chart every minute (60000 ms) - DISABLED, server handles minute transitions
			// setInterval(() => {
			// 	fetch('/api/stackedbarchart/advance?datastar=true')
			// 		.then(response => {
			// 			if (!response.ok) console.error('Advance minute failed');
			// 		})
			// 		.catch(err => console.error('Advance minute error:', err));
			// }, 60000);
		</script>
	</head>
	<body>
		<div class="container">
			<h1>Stacked Bar Chart Component</h1>
			<p class="subtitle">Washing Machine Delay Monitor - Last 10 Minutes</p>
			
			<div class="controls">
				<button data-on:click="@get('/api/stackedbarchart/randomize')">
					Randomize Data
				</button>
			</div>
			
			<div class="chart-container">
				` + data.HTML + `
			</div>
			
			<a href="/" class="back-link">← Back to Component Library</a>
		</div>
	</body>
	</html>`))
}

// piechartDemoHandler serves the original pie chart demo page (for backward compatibility)
func piechartDemoHandler(w http.ResponseWriter, r *http.Request) {
	piechartComp := piechart.New()
	data := piechartComp.GenerateRandomData()

	component := templates.Page(data)
	w.Header().Set("Content-Type", "text/html")
	component.Render(r.Context(), w)
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
	r := chi.NewRouter()

	// Add middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(requestTimer)

	// Setup component API routes
	setupRouter(r)

	// Home page with component library
	r.Get("/", homeHandler)

	// Component pages (standalone views)
	r.Get("/component/piechart", piechartPageHandler)
	r.Get("/component/barchart", barchartPageHandler)
	r.Get("/component/button", buttonPageHandler)
	r.Get("/component/form", formPageHandler)
	r.Get("/component/targetbarchart", targetbarchartPageHandler)
	r.Get("/component/stackedbarchart", stackedbarchartPageHandler)

	// Original pie chart demo (backward compatibility)
	r.Get("/demo", piechartDemoHandler)

	// Serve global static files
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	port := ":8080"
	fmt.Printf("Server starting on http://localhost%s\n", port)
	fmt.Println("Press Ctrl+C to stop")

	log.Fatal(http.ListenAndServe(port, r))
}
