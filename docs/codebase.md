# Codebase Overview

This document provides a high-level overview of the StatusApp codebase.

The application is primarily written in Go and uses various libraries for TUI (Text User Interface) development, inter-process communication, and system interactions.

## Key Components:

- **`cmd/main.go`**: The entry point of the application. Initializes the TUI, loads configurations, and starts the main application loop.
- **`internal/renderers/`**: Contains logic for rendering different status modules (e.g., clock, schedule, weather, Tailscale status) within the TUI. Each renderer is responsible for fetching data and presenting it in a formatted way.
- **`internal/models/`**: Defines the data structures used throughout the application, such as configuration settings, weather data, schedule entries, and Tailscale status.
- **`internal/schedule/`**: Manages the loading and parsing of user-defined schedules.
- **`internal/tailscale/`**: Handles interactions with the Tailscale API to retrieve network status information.
- **`internal/weather/`**: Manages interactions with a weather API to fetch current weather conditions.
- **`configs/constants.go`**: Stores application-wide constants and default values.
- **`assets/`**: Contains static assets like FIGlet fonts (`big.flf`) and other configuration files (`mobius.txt`, `weather.json`).

## Dependencies:

The project utilizes several external Go modules, managed via `go.mod` and `go.sum`, primarily for:
- **`charmbracelet/bubbletea`**: A powerful library for building terminal-based user interfaces.
- **`charmbracelet/lipgloss`**: A style definition library for color and formatting in the terminal.
- Other utilities for system notifications, environment variable loading, and more.

## Architecture:

The application follows a modular architecture where different concerns are separated into distinct packages. Data flow generally involves:
1. Configuration loading at startup.
2. Periodic data fetching by renderer-specific or service-specific modules (`internal/tailscale`, `internal/weather`).
3. Data processing and formatting by `internal/renderers`.
4. Rendering of the formatted data to the terminal using `bubbletea` and `lipgloss`.

For more details on project structure, refer to `docs/structure.md`.
