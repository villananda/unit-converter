# Unit Converter

A small web-based unit converter built with Go and Gin, with a dynamic frontend powered by HTMX. It converts length, weight, and temperature units through a tabbed interface, validates input on the server, and returns results without a full page reload.

## Tech Stack

| Layer | Technology |
| --- | --- |
| Language | Go 1.27.0 |
| Web framework | [Gin](https://github.com/gin-gonic/gin) v1.12.0 |
| Templating | Go `html/template` (loaded via `gin.LoadHTMLGlob`) |
| Frontend | HTMX 2.0.4 |
| Styling | Plain CSS (inline in the template) |
| Testing | Go's built-in `testing` package |

## Features

- **Length**: millimeter, centimeter, meter, kilometer, inch, foot, yard, mile
- **Weight**: milligram, gram, kilogram, ounce, pound
- **Temperature**: Celsius, Fahrenheit, Kelvin
- Tabbed UI with deep-linkable category URLs (`/?tab=length`)
- Server-side validation of numeric input (rejects empty, non-numeric, `NaN`, and `Inf` values)
- Rejects physical impossibilities such as temperatures below absolute zero
- HTMX-powered partial updates — no full page reloads

## Prerequisites

- [Go](https://go.dev/dl/) 1.27.0 or newer
- A modern browser

## How to Run

1. Clone the repository and move into the project directory:

   ```bash
   git clone <repository-url>
   cd unit-converter
   ```

2. Download the dependencies:

   ```bash
   go mod download
   ```

3. Start the server:

   ```bash
   go run .
   ```

4. Open [http://localhost:8080](http://localhost:8080) in your browser.

The server listens on port `8080`. To use a different port, change the address passed to `r.Run` in `main.go`.

## Project Structure

```
unit-converter/
├── main.go            # HTTP server, routing, request handling, view models
├── converter.go       # Unit/category definitions and conversion logic
├── converter_test.go  # Unit tests for the conversion logic
├── templates/
│   └── index.html     # Tabbed UI (Go template + HTMX + CSS)
├── go.mod
├── go.sum
└── README.md
```

### How conversion works

Every unit defines `toBase` and `fromBase` functions. A conversion goes `source → base unit → target`, which keeps the logic simple and easy to extend:

- Length and weight use a linear factor (base units: meter and gram)
- Temperature uses Kelvin as the base, so the non-negative check naturally means "not below absolute zero"

To add a unit, append a `linear(...)` entry to a category's `Units` slice in `converter.go` (or a custom `Unit` for non-linear conversions). To add a category, append a new `Category` to `categories`.

## Routing

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/` | Renders the page; optional `?tab=<category>` selects the active tab |
| `POST` | `/convert/:category` | Converts the submitted form and re-renders the page |

Form fields for `POST /convert/:category`: `value`, `from`, `to`.

## Testing

Run the test suite:

```bash
go test ./...
```

Run with verbose output:

```bash
go test -v ./...
```

## Reference

This project is based on the [Unit Converter](https://roadmap.sh/projects/unit-converter) project idea from [roadmap.sh](https://roadmap.sh).
