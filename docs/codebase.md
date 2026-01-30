# Codebase Overview

This document provides a high-level overview of the StatusApp codebase.

The application is primarily written in Go and uses various libraries for TUI (Text User Interface) development, inter-process communication, and system interactions.

## Key Components:

- **`cmd/main.go`**: The entry point of the application. Initializes the TUI, loads configurations, and starts the main application loop.
- **`cmd/view.go`**: Contains the main view logic for the application, orchestrating the layout of different components.
- **`cmd/fetch.go`**: Contains the logic for fetching data from various sources.
- **`internal/clock/`**: Contains logic for rendering the digital clock.
- **`internal/common/`**: Contains common utility functions.
- **`internal/hosthatch/`**: Handles interactions with the HostHatch API to retrieve server status information.
- **`internal/schedule/`**: Manages the loading and parsing of user-defined schedules.
- **`internal/tailscale/`**: Handles interactions with the Tailscale API to retrieve network status information.
- **`internal/truenas/`**: Handles interactions with the TrueNAS API to retrieve application status.
- **`internal/upcloud/`**: Handles interactions with the UpCloud API to retrieve server status information.
- **`internal/weather/`**: Manages interactions with a weather API to fetch current weather conditions.
- **`configs/constants.go`**: Stores application-wide constants and default values.
- **`assets/`**: Contains static assets like FIGlet fonts (`big.flf`) and other configuration files (`weather.json`, `weather_yr.csv`).

## Dependencies:

The project utilizes several external Go modules, managed via `go.mod` and `go.sum`, primarily for:
- **`charmbracelet/bubbletea`**: A powerful library for building terminal-based user interfaces.
- **`charmbracelet/lipgloss`**: A style definition library for color and formatting in the terminal.
- Other utilities for system notifications, environment variable loading, and more.

## Architecture:

The application follows a modular architecture where different concerns are separated into distinct packages. Data flow generally involves:
1. Configuration loading at startup.
2. Periodic data fetching by client code within each module (e.g., `internal/tailscale/client.go`, `internal/weather/client.go`).
3. Data processing and formatting by view functions within each module (e.g., `internal/tailscale/view.go`).
4. Rendering of the formatted data to the terminal using `bubbletea` and `lipgloss`, orchestrated by `cmd/view.go`.

For more details on project structure, refer to `docs/structure.md`.
