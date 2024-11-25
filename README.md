# go_otel_auto_instrument

`go_otel_auto_instrument` let's you automatically instrument your Go code with Opentelemtry.

## Installation

### 1. Install the `go-otel-auto-instrument` executable

To use `go-otel-auto-instrument` in your build process, install it as a binary executable:

```bash
go install github.com/pijng/go_otel_auto_instrument/cmd/go-otel-auto-instrument@latest
```

### 2. Add `go_otel_auto_instrument` as a dependency

To use the `go_otel_auto_instrument` library in your Go project, add it via Go modules:

```bash
go get github.com/pijng/go_otel_auto_instrument
go mod tidy
```

## Usage

### Add blank import to main package

```go
package main

import (
  _ "github.com/pijng/go_otel_auto_instrument"
)
```

### Building with `go-otel-auto-instrument`

To automatically instrument your Go code run the command:

```bash
go build -a -toolexec="go-otel-auto-instrument <absolute/path/to/project>" main.go
```

**Important:**
  * `-a` flag is required to recompile all your project, otherwise go compiler might do nothing and use cached build
  * `<absolute/path/to/project>` is an absolute path to the root of your project. If you run `go build` from the root of the project – simply specify `$PWD` as an argument.

### Run the final binary and specify OTEL_* env variables:

```bash
OTEL_SERVICE_NAME="My Service" OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318 ./main
```
