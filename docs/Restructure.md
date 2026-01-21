# Project Restructuring Recommendations

## 1. Introduction

This document provides an analysis of the current project structure and offers recommendations for reorganization. The goal is to align the project with common Go best practices, including principles from the `golang-standards/project-layout`, to improve modularity, cohesion, and maintainability.

The current structure is functional, but as the application grows, some of the patterns may introduce friction.

## 2. Analysis and Recommendations

### On `internal/models`

**Observation:** The project centralizes all data structures in a single `internal/models` package. While this can seem convenient, it often leads to a low-cohesion package that becomes a bottleneck for changes and can create complex dependency chains or even import cycles (as noted in `migration.md`).

**Recommendation: Co-locate Models within Domain Packages.**

Models should be defined by the package that is most concerned with them. This practice, known as co-location, increases the cohesion of your packages.

- **Before:**
  ```
  internal/
  ├── models/
  │   └── tailscale.go  // Defines TailscaleStatus struct
  └── tailscale/
      └── webrequests.go  // Imports models, fetches data
  ```

- **After:**
  ```
  internal/
  └── tailscale/
      ├── types.go        // Defines TailscaleStatus struct
      └── client.go       // Fetches data and uses TailscaleStatus
  ```

This makes the `internal/tailscale` package a self-contained unit. The `internal/models` package would be removed entirely, with its contents distributed to the relevant domain packages.

### On `internal/renderers`

**Observation:** The Bubble Tea view/render functions are collected in `internal/renderers`. This separates presentation logic from business logic, which is a good instinct. However, `renderers/tailscale.go` is only ever concerned with `tailscale` data.

**Recommendation: Integrate Renderers into Domain Packages.**

Just like models, the rendering logic is tightly coupled to the domain it represents. Moving it inside the domain package further improves cohesion.

- **Before:**
  ```
  internal/
  ├── renderers/
  │   └── tailscale.go  // Contains RenderTailscale()
  └── tailscale/
      └── webrequests.go
  ```

- **After:**
  ```
  internal/
  └── tailscale/
      ├── client.go
      ├── types.go
      └── view.go         // Contains View() or Render() for tailscale
  ```
Your main application can then call `tailscale.View(model)` instead of `renderers.RenderTailscale(model)`.

### On "Clients" vs. Current Organization

**Observation:** You asked whether features like `Tailscale` should be split into `Clients`. Currently, packages like `internal/tailscale` contain a `webrequests.go` file, which acts as a client for an external API. We've decided to proceed with a direct method call pattern (e.g., `client.GetMachines()`) rather than a fluent API style (`client.Machines().Get()`) due to the API's relative simplicity.

**Recommendation: Formalize the Client Role with a Struct and Direct Methods.**

Your intuition is correct. The code that interacts with external services is a "client". Instead of creating a top-level `internal/clients` directory, you can adopt a clean internal structure for each domain package. This approach is highly cohesive and scales well.

A `Client` struct should be defined, holding dependencies like the HTTP client, API keys, and base URL. A constructor function (e.g., `NewClient`) should be used to initialize this struct. Methods performing API calls (e.g., `GetMachines`) should then be defined directly on this `Client` struct.

**Example Client Structure (`internal/tailscale/client.go`):**

```go
package tailscale

import (
	"context"
	"net/http"
	"time"
)

// Client is a client for the Tailscale API.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	// ... other dependencies
}

// NewClient creates a new, configured Tailscale client.
func NewClient(apiKey, tailnet string) *Client {
	return &Client{
		baseURL: "https://api.tailscale.com/api/v2", // Example base URL
		apiKey:  apiKey,
		// ... initialize httpClient with appropriate timeouts, etc.
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetMachines fetches the list of devices from the Tailscale API.
// This is a direct method call on the client instance.
func (c *Client) GetMachines(ctx context.Context) ([]Machine, error) {
	// Implementation to make HTTP request, handle authentication,
	// and unmarshal response into []Machine (defined in types.go)
	// Example:
	// url := fmt.Sprintf("%s/tailnet/%s/devices", c.baseURL, c.tailnet)
	// req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	// req.SetBasicAuth(c.apiKey, "")
	// resp, err := c.httpClient.Do(req)
	// ...
	return nil, nil // Placeholder
}
```

This makes each domain package a fully-featured, self-contained component responsible for its own data, logic, external communication, and presentation.

### Interfaces for Testability and Decoupling

To further improve testability and decouple your application logic from the concrete client implementations, it's highly recommended to define **interfaces** for your clients.

**Why use interfaces?**

*   **Mocking in Tests:** During unit tests, you can provide a "mock" implementation of the interface that returns predefined data or errors, without needing to make actual network calls. This makes tests faster, more reliable, and deterministic.
*   **Decoupling:** Your business logic (e.g., `main.go` or other service layers) will depend on the interface, not the concrete `*Client` struct. This means if you ever change the underlying client implementation, as long as it satisfies the interface, your business logic doesn't need to change.

**Example (`internal/tailscale/client.go`):**

```go
package tailscale

import (
	"context"
)

// MachineGetter defines the contract for fetching machines.
// Any type that implements this method can be used where a MachineGetter is expected.
type MachineGetter interface {
	GetMachines(ctx context.Context) ([]Machine, error)
	// Add other methods here as the client grows, e.g.,
	// GetMachine(ctx context.Context, id string) (Machine, error)
}

// The Client struct (defined above) implicitly implements the MachineGetter interface
// because it has a GetMachines method with the correct signature.

// Example of how business logic would use the interface:
/*
func ProcessTailscaleData(ctx context.Context, mg MachineGetter) error {
	machines, err := mg.GetMachines(ctx)
	if err != nil {
		return fmt.Errorf("failed to get machines: %w", err)
	}
	// ... process machines ...
	return nil
}
*/
```

### On `pkg/` vs `internal/`

The `golang-standards/project-layout` guide distinguishes between `pkg/` and `internal/`.

- **`internal/`**: For private application code. It's a Go compiler-enforced rule that you cannot import code from another project's `internal` directory.
- **`pkg/`**: For library code that is safe to be used by external applications.

You mentioned, "There is nothing in this project that should be shared with other projects." Therefore, **your current use of `internal/` for all your application logic is absolutely correct** and follows best practices.

## 3. Adding Tests to the Project

Since the project currently has no tests, incorporating them is a high-priority step that will significantly improve reliability and make future development safer. Go has excellent built-in support for testing.

### Go Testing Basics

-   Tests for a file `source.go` are placed in a corresponding file named `source_test.go`.
-   Test files are part of the same package as the code they are testing.
-   Run tests for the entire project from the root directory with `go test ./...`.

### a) Testing the API Client (e.g., `internal/tailscale/client_test.go`)

To test the client that makes real HTTP calls, we should not hit the actual Tailscale API in our tests. Instead, we can use the `net/http/httptest` package to spin up a temporary, local HTTP server that returns a predictable response.

This approach tests everything in our client (request creation, authentication, JSON decoding) without network dependency.

**Example (`internal/tailscale/client_test.go`):**
```go
package tailscale

import (
    "context"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestGetMachines_Success(t *testing.T) {
    // 1. Create a mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Assert that the request is what we expect
        if r.URL.Path != "/api/v2/tailnet/my-tailnet/devices" {
            t.Fatalf("unexpected path: %s", r.URL.Path)
        }
        // Respond with a canned JSON payload
        fmt.Fprintln(w, `{"devices": [{"id": "1", "name": "device-1"}, {"id": "2", "name": "device-2"}]}`)
    }))
    defer server.Close()

    // 2. Create a client instance pointing to our mock server
    client := &Client{
        baseURL:    server.URL, // Use the test server's URL
        apiKey:     "test-key",
        tailnet:    "my-tailnet",
        httpClient: server.Client(),
    }

    // 3. Call the method being tested
    machines, err := client.GetMachines(context.Background())
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }

    // 4. Assert the results are correct
    if len(machines) != 2 {
        t.Fatalf("expected 2 machines, got %d", len(machines))
    }
    if machines[0].Name != "device-1" {
        t.Errorf("expected machine name 'device-1', got '%s'", machines[0].Name)
    }
}
```

### b) Testing Business Logic with Mock Clients

This is where the interface we defined earlier (`MachineGetter`) becomes powerful. For business logic that *uses* the client, we don't need a mock HTTP server. We just need a "mock client" struct that implements the interface.

**Example Mock and Test (`internal/applogic/processor_test.go`):**
```go
package applogic // A hypothetical package for your business logic

import (
    "context"
    "errors"
    "testing"
    "github.com/your-repo/internal/tailscale" // Import your tailscale package
)

// mockMachineGetter is a mock implementation of the tailscale.MachineGetter interface.
type mockMachineGetter struct {
    machines []tailscale.Machine
    err      error
}

// GetMachines satisfies the interface. It returns the canned data/error.
func (m *mockMachineGetter) GetMachines(ctx context.Context) ([]tailscale.Machine, error) {
    return m.machines, m.err
}

// Test function for business logic that depends on the interface
func TestProcessTailscaleData(t *testing.T) {
    t.Run("success case", func(t *testing.T) {
        // 1. Setup the mock to return successful data
        mock := &mockMachineGetter{
            machines: []tailscale.Machine{{Name: "online-device"}},
        }
        
        // 2. Call the function with the mock
        summary, err := ProcessTailscaleData(context.Background(), mock)
        if err != nil {
            t.Fatalf("expected no error, got %v", err)
        }

        // 3. Assert the business logic worked correctly
        if summary.Online != 1 {
            t.Errorf("expected 1 online device, got %d", summary.Online)
        }
    })

    t.Run("error case", func(t *testing.T) {
        // 1. Setup the mock to return an error
        mock := &mockMachineGetter{
            err: errors.New("API unavailable"),
        }

        // 2. Call the function with the mock
        _, err := ProcessTailscaleData(context.Background(), mock)

        // 3. Assert that the error was handled correctly
        if err == nil {
            t.Fatal("expected an error, but got none")
        }
    })
}
```

### c) Testing the Renderer/View (`internal/tailscale/view_test.go`)

TUI rendering logic can be tested by treating the view functions as pure functions: they take state (`model`) and return a `string`. Your tests can check if the output string contains the expected text for a given state.

**Example (`internal/tailscale/view_test.go`):**
```go
package tailscale

import (
    "strings"
    "testing"
    // You would import your main model definition here
)

func TestView(t *testing.T) {
    // 1. Create a model with some state to render
    // This model would be the one used by your Bubble Tea application.
    // The exact structure will depend on your implementation.
    model := AppModel{
        tailscaleData: []Machine{
            {Name: "server-1", OS: "linux"},
            {Name: "laptop", OS: "darwin"},
        },
        error: nil,
    }

    // 2. Call the view function
    output := View(model) // Assuming a View function exists in the package

    // 3. Assert that the output contains the expected content
    if !strings.Contains(output, "server-1 (linux)") {
        t.Errorf("view output does not contain 'server-1 (linux)'")
    }
    if !strings.Contains(output, "laptop (darwin)") {
        t.Errorf("view output does not contain 'laptop (darwin)'")
    }
    if strings.Contains(output, "Error:") {
        t.Errorf("view output should not contain an error message")
    }
}
```

## 4. Proposed Future Project Structure

Based on the recommendations above, the reorganized `internal` directory would look like this, now including test files:

```
/home/manfred/Development/StatusApp/
├───.gitignore
├───go.mod
├───go.sum
...
├───cmd/
│   └───main.go
├───configs/
│   └───constants.go
...
├───internal/
│   ├───common/
│   │   ├───time.go
│   │   └───time_test.go
│   ├───schedule/
│   │   ├───loader.go
│   │   ├───loader_test.go
│   │   └───types.go
│   ├───tailscale/
│   │   ├───client.go
│   │   ├───client_test.go
│   │   ├───types.go
│   │   ├───view.go
│   │   └───view_test.go
│   ├───truenas/
│   │   ├───apps.go
│   │   ├───apps_test.go
│   │   └───types.go
│   ├───weather/
│   │   ├───client.go
│   │   ├───client_test.go
│   │   ├───watertemp.go
│   │   ├───watertemp_test.go
│   │   ├───types.go
│   │   ├───view.go
│   │   └───view_test.go
...
```

## 5. Conclusion

Adopting this structure will make your project more modular and easier to maintain. Each domain package will become a self-sufficient component, which simplifies testing, debugging, and future development. The monolithic `models` and `renderers` packages, which can become sources of complexity, are eliminated in favor of a more cohesive, domain-driven design. The use of interfaces and the addition of a comprehensive test suite will ensure the application is robust and reliable.
