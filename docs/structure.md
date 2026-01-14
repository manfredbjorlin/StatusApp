# Project Structure

This document outlines the directory and file structure of the StatusApp project.

## Root Directory

- **`.gitignore`**: Specifies intentionally untracked files to ignore.
- **`go.mod`**, **`go.sum`**: Go module definition files, managing dependencies.
- **`LICENSE`**: The project's license file.
- **`README.md`**: The main project documentation, providing an overview and instructions.
- **`assets/`**: Contains static assets used by the application.
    - **`big.flf`**: FIGlet font file.
    - **`weather.json`**: Weather-related configuration or data.
- **`cmd/`**: Contains the main application entry points.
    - **`main.go`**: The primary executable file that starts the StatusApp.
- **`configs/`**: Stores application configuration files.
    - **`constants.go`**: Defines various constants used throughout the application.
- **`deployments/`**: Contains deployment-related files.
    - **`.env`**: Environment variables for deployment.
    - **`StatusApp`**: Deployment script or configuration for the application.
- **`docs/`**: Project documentation files.
    - **`codebase.md`**: Overview of the codebase (this file).
    - **`migration.md`**: Documentation related to data migrations or upgrades.
    - **`structure.md`**: Describes the project's directory structure (this file).
- **`internal/`**: Contains internal application logic, not intended for external consumption.
    - **`models/`**: Defines data structures and models.
        - **`models.go`**: Go struct definitions for various data entities.
    - **`renderers/`**: Logic for rendering different status components.
        - **`clock.go`**: Handles rendering of the clock component.
        - **`schedule.go`**: Handles rendering of the schedule component.
        - **`tailscale.go`**: Handles rendering of the Tailscale status component.
        - **`weather.go`**: Handles rendering of the weather component.
    - **`schedule/`**: Logic related to scheduling.
        - **`loader.go`**: Loads and manages schedule data.
    - **`tailscale/`**: Interacts with the Tailscale API.
        - **`webrequests.go`**: Handles web requests to the Tailscale API.
    - **`weather/`**: Interacts with a weather API.
        - **`webrequests.go`**: Handles web requests to the weather API.
- **`scripts/`**: Contains various utility scripts.
    - **`Build.sh`**: Script to build the application.
    - **`BuildAndRun.sh`**: Script to build and run the application.
    - **`BuildToWindows.sh`**: Script to build the application for Windows.
- **`vendor/`**: Go module dependencies. This directory is managed by `go modules` and contains vendored copies of external libraries.

## `internal/` Directory Breakdown

The `internal/` directory is central to the application's logic, organizing code by domain or feature. Each subdirectory within `internal/` typically encapsulates a specific part of the application's functionality, promoting modularity and maintainability.

- **`models/`**: Standard Go structs that represent the data used and manipulated by the application.
- **`renderers/`**: Modules responsible for taking application data and transforming it into a format suitable for display in the terminal UI. Each `.go` file here corresponds to a distinct UI component or "widget."
- **`schedule/`**: Contains the business logic for managing events, parsing schedule files, and determining what to display.
- **`tailscale/`**: Abstraction layer for interacting with Tailscale's API, handling authentication and data retrieval.
- **`weather/`**: Abstraction layer for interacting with a weather service API, handling requests and parsing responses.

This structure aims to keep the codebase organized, making it easier to navigate, understand, and maintain.
