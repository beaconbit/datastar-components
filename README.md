# Pie Chart Demo with Go, Templ, Chi, and Datastar

A demonstration project showing:
- Go backend with random pie chart generation
- Chi router for HTTP routing with middleware
- Templ templates for type-safe HTML
- Datastar for hypermedia-driven interactivity
- SVG pie chart with randomized data

## Features

- Interactive pie chart rendered as SVG
- "Randomize" button using Datastar's `data-on:click="@get('/randomize')"`
- Backend returns full HTML snippet for Datastar to morph into DOM
- Responsive design with CSS styling
- Legend showing sector colors and percentages
- Chi router with middleware (logging, compression, recovery)

## Project Structure

- `main.go` - HTTP server with Chi router and handlers
- `components/` - Pie chart components and helpers
- `templates/` - Page templates
- `go.mod` - Go dependencies

## Running the Project

```bash
# Install dependencies
go mod tidy

# Generate templ files
templ generate

# Run the server
go run main.go

# Open http://localhost:8080
```

## How It Works

1. The home page (`/`) renders a pie chart with random data using Templ
2. The "Randomize Chart" button has `data-on:click="@get('/randomize')"`
3. Clicking sends a GET request to `/randomize` with Datastar signals
4. Server generates new random data and returns `datastar-patch-elements` SSE event
5. Datastar morphs the new chart HTML into the DOM (matching `id="pie-chart"`)
6. No JavaScript framework needed - just HTML over the wire

## Technologies

- **Go**: Backend server and data generation
- **Chi**: Lightweight, idiomatic HTTP router with middleware
- **Templ**: Type-safe HTML templates
- **Datastar**: Hypermedia framework for interactivity
- **SVG**: Scalable vector graphics for pie chart

## Endpoints

- `GET /` - Main page with pie chart
- `GET /randomize` - Returns new pie chart HTML via SSE
- `GET /signals` - Example signals endpoint

## Middleware

- `RequestID` - Adds request ID to each request
- `RealIP` - Handles X-Forwarded-For header
- `Logger` - HTTP request logging
- `Recoverer` - Recovers from panics
- `Compress` - Gzip compression
- `requestTimer` - Custom middleware logging request duration