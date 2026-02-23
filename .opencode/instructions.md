# OpenCode Configuration
# This file contains documentation for Go, Templ, and Datastar for use with OpenCode

docs:
  golang:
    title: "Go (Golang) Documentation"
    summary: |
      Go is a statically typed, compiled programming language designed at Google.
      It is known for its simplicity, concurrency support, and efficient performance.
    key_concepts:
      - packages: "Every Go program is made up of packages. Programs start running in package main."
      - imports: "Import statements declare package dependencies."
      - exported_names: "Names starting with capital letters are exported (public)."
      - functions: "Functions can return multiple values. Use `func` keyword."
      - variables: "Declare with `var` or using short assignment `:=`."
      - types: "Basic types: bool, string, int, float64, etc. Structs and interfaces."
      - concurrency: "Goroutines (go keyword) and channels for communication."
      - error_handling: "Errors are values; typically returned as second return value."
      - defer: "Schedule a function call to be run after the function completes."
    common_patterns:
      - "http.HandlerFunc for HTTP handlers"
      - "context.Context for request-scoped values and cancellation"
      - "error wrapping with fmt.Errorf and errors.Is/As"
      - "JSON marshaling/unmarshaling with encoding/json"
      - "testing with testing package and table-driven tests"
    resources:
      - "Official tour: https://go.dev/tour"
      - "Effective Go: https://go.dev/doc/effective_go"
      - "Go by Example: https://gobyexample.com"

  templ:
    title: "Templ Templating Syntax"
    summary: |
      Templ is a Go HTML templating language that compiles to Go code.
      It provides type-safe HTML templates with component composition.
    syntax:
      - "Components: `templ ComponentName() { ... }`"
      - "HTML elements: `div { ... }` (no angle brackets)"
      - "Text content: `{ variable }` or `{ expression }`"
      - "Attributes: `class("my-class")` or `attr("name", value)`"
      - "Conditionals: `if condition { ... } else { ... }`"
      - "Loops: `for _, item := range items { ... }`"
      - "Children: `@children` slot for component composition"
      - "CSS: `style` function with type-safe CSS"
      - "Scripts: `script` function for inline JavaScript"
      - "Comments: `//` for single-line, `/* */` for multi-line"
    example: |
      // Component definition
      package main

      import "github.com/a-h/templ"

      templ Hello(name string) {
          div(class("greeting")) {
              h1 { "Hello, " { name } }
          }
      }

      // Usage in handler
      func handler(w http.ResponseWriter, r *http.Request) {
          component := Hello("World")
          component.Render(r.Context(), w)
      }
    resources:
      - "Official docs: https://templ.guide"
      - "GitHub: https://github.com/a-h/templ"

  datastar:
    title: "Datastar Documentation"
    summary: |
      Datastar is a hypermedia-first framework for building backend-driven,
      interactive UIs using declarative `data-*` HTML attributes.
      It provides frontend reactivity and backend reactivity in a lightweight package.
    core_concepts:
      - "Hypermedia-first: Backend drives frontend state via HTML patches"
      - "Signals: Reactive variables denoted with `$` prefix"
      - "Data attributes: Declarative HTML attributes for reactivity"
      - "SSE events: Server-Sent Events for real-time updates"
      - "Patch Elements: Backend can patch HTML elements into DOM"
      - "Patch Signals: Backend can update frontend signals"
    key_attributes:
      - "data-bind: Two-way data binding for input elements"
      - "data-text: Sets text content to signal value"
      - "data-computed: Creates derived read-only signal"
      - "data-show: Shows/hides element based on expression"
      - "data-class: Adds/removes CSS class based on expression"
      - "data-attr: Binds HTML attribute to expression"
      - "data-signals: Patches signals into frontend state"
      - "data-on: Attaches event listener with expression"
      - "data-indicator: Signal for request loading state"
    actions:
      - "@get(): Sends GET request to backend"
      - "@post(): Sends POST request"
      - "@put(), @patch(), @delete(): Other HTTP methods"
      - "@toggleAll(): Toggles multiple signals"
    expressions:
      - "JavaScript-like syntax with signal references (`$signal`)"
      - "Access element via `el` variable"
      - "Multiple statements separated by semicolons"
      - "Ternary operator `?:`, logical OR `||`, AND `&&`"
    installation: |
      <!-- CDN -->
      <script type="module" src="https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js"></script>

      <!-- Package manager -->
      import 'https://cdn.jsdelivr.net/gh/starfederation/datastar@1.0.0-RC.7/bundles/datastar.js'
    backend_integration:
      - "SDKs available for Go, C#, Java, Kotlin, PHP, Python, Ruby, Rust, JavaScript, Clojure"
      - "SSE events: datastar-patch-elements, datastar-patch-signals"
      - "Read signals from request, send patches via SSE"
    example: |
      <div data-signals:hal="'...'">
          <button data-on:click="@get('/endpoint')">
              HAL, do you read me?
          </button>
          <div data-text="$hal"></div>
      </div>

      // Backend (Go SDK)
      sse := datastar.NewSSE(w, r)
      sse.PatchSignals([]byte(`{hal: 'Affirmative, Dave. I read you.'}`))
    source: "Documentation fetched from https://data-star.dev/docs.md"
    resources:
      - "Official docs: https://data-star.dev/docs"
      - "Reference: https://data-star.dev/reference"
      - "Examples: https://data-star.dev/examples"
      - "VSCode extension for autocompletion"

# OpenCode specific configuration
opencode:
  # Project type detection
  project_types:
    - "go"
    - "templ"
    - "datastar"
  # Preferred documentation sources
  doc_priority:
    - "local"
    - "official"
    - "community"
  # Auto-suggest patterns
  suggestions:
    go: "Use gofmt and go vet for code quality"
    templ: "Run templ generate after template changes"
    datastar: "Use data-* attributes for reactivity, keep state in backend"
