# livelog

[![Go Reference](https://pkg.go.dev/badge/github.com/bmvkrd/livelog.svg)](https://pkg.go.dev/github.com/bmvkrd/livelog)

`livelog` is a Go library for terminal output that keeps a **live zone** pinned at the bottom of the screen while scrolling log lines accumulate above it. It is useful for CLI tools that display progress bars, spinners, or status lines alongside regular output.

When the output is not a TTY (e.g. piped to a file or CI log), all output degrades gracefully to plain newline-separated text with no ANSI codes.

## Features

- **Live zone** — one or more lines pinned at the bottom, redrawn in-place on every update
- **Scrolling log** — `Log`/`Logf` lines accumulate above the live zone without disturbing it
- **`io.Writer` adapter** — wrap a `Display` in a `Writer` to route `fmt.Fprintf` / `log` output through the log zone
- **`ProgressBar`** — simple text progress bar that renders to a string (use it anywhere)
- **ANSI utilities** — `StripANSI`, `VisibleWidth`, and `Truncate` for correct handling of colored strings
- **Concurrency-safe** — all `Display` methods are safe for concurrent use

## Installation

```sh
go get github.com/bmvkrd/livelog@latest
```

## Quick start

```go
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/bmvkrd/livelog"
)

func main() {
    d := livelog.New(os.Stdout)
    defer d.Flush() // always call Flush to clear the live zone on exit

    pb := &livelog.ProgressBar{Total: 5, Width: 30}

    for i := 0; i < 5; i++ {
        pb.Current = float64(i)
        d.SetLive([]string{
            fmt.Sprintf("Processing item %d/5", i+1),
            pb.String(),
        })
        time.Sleep(500 * time.Millisecond)
        d.Logf("✓ item %d done", i+1)
    }
}
```

Run the bundled demo to see it in action:

```sh
go run ./examples/demo
```

## API

### Display

`Display` is the central type. Create one with `New` and pass it your output file (usually `os.Stdout`).

```go
d := livelog.New(os.Stdout)
```

#### Options

| Option | Description |
|--------|-------------|
| `WithForceTTY(bool)` | Override automatic TTY detection. Useful in tests. |

#### Methods

| Method | Description |
|--------|-------------|
| `Log(args ...any)` | Write a line to the scrolling log zone. |
| `Logf(format string, args ...any)` | Write a formatted line to the scrolling log zone. |
| `SetLive(lines []string)` | Replace the live zone with the given lines. Pass `nil` to clear. |
| `ClearLive()` | Erase the live zone immediately. |
| `Flush()` | Alias for `ClearLive`; call before program exit. |
| `Writer() *displayWriter` | Return an `io.Writer` that routes complete lines through `Log`. |
| `IsTTY() bool` | Report whether the output is a terminal. |
| `TerminalWidth() int` | Return the terminal width in columns (fallback: 80). |

### ProgressBar

`ProgressBar` renders a text-based progress bar. It implements `fmt.Stringer` and can be used with any output method.

```go
pb := &livelog.ProgressBar{
    Total:      100,
    Current:    45,
    Width:      40,       // optional, default 40
    FilledChar: "█",      // optional, default "█"
    EmptyChar:  "░",      // optional, default "░"
}
d.SetLive([]string{pb.String()})
// ██████████████████░░░░░░░░░░░░░░░░░░░░  45%  45s / 100s
```

| Method | Description |
|--------|-------------|
| `SetRatio(ratio float64)` | Set `Current` from a 0.0–1.0 ratio (clamped). |
| `String() string` | Render the bar as a string. |

### Writer adapter

Use `d.Writer()` to get an `io.Writer` that feeds into the log zone. This lets you pass a `Display` to the standard `log` package or to third-party libraries that accept an `io.Writer`.

```go
d := livelog.New(os.Stdout)
logger := log.New(d.Writer(), "", log.LstdFlags)
logger.Println("this appears in the log zone")
```

### ANSI utilities

These are exported as part of the `livelog` package and can be used independently.

| Function | Description |
|----------|-------------|
| `StripANSI(s string) string` | Remove all ANSI escape sequences from a string. |
| `VisibleWidth(s string) int` | Column width of a string after stripping ANSI codes. |
| `Truncate(s string, maxWidth int) string` | Shorten a string (preserving ANSI sequences) to fit within `maxWidth` visible columns. Appends `…` and a reset sequence when truncated. |

## Non-TTY behaviour

When the output is not a terminal, `Display` writes plain text with no ANSI codes:

- `Log` / `Logf` write the line followed by `\n`.
- `SetLive` prints each line followed by `\n` (no in-place redraw).
- `ClearLive` / `Flush` are no-ops.

This makes `livelog` safe to use in scripts, Docker containers, and CI pipelines without any special handling.

## Versioning

This module follows [semantic versioning](https://semver.org). The current version constant is available at runtime:

```go
fmt.Println(livelog.Version) // "0.1.0"
```

## License

MIT
