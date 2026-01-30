# Project Structure

This document outlines the directory and file structure of the StatusApp project.

## Root Directory

- **`.gitignore`**: Specifies intentionally untracked files to ignore.
- **`go.mod`**, **`go.sum`**: Go module definition files, managing dependencies.
- **`LICENSE`**: The project's license file.
- **`README.md`**: The main project documentation, providing an overview and instructions.
- **`assets/`**: Contains static assets used by the application.
    - **`big.flf`**: FIGlet font file.
    - **`weather_yr.csv`**: Weather data for YR service.
    - **`weather.json`**: Weather-related configuration or data.
- **`cmd/`**: Contains the main application entry points and commands.
    - **`main.go`**: The primary executable file that starts the StatusApp.
    - **`fetch.go`**: Logic for fetching data.
    - **`update.go`**: Logic for updating data.
    - **`view.go`**: Logic for displaying data.
    - **`types.go`**: Type definitions for the `cmd` package.
- **`configs/`**: Stores application configuration files.
    - **`constants.go`**: Defines various constants used throughout the application.
- **`deployments/`**: Contains deployment-related files (currently empty).
- **`docs/`**: Project documentation files.
    - **`codebase.md`**: Overview of the codebase.
    - **`migration.md`**: Documentation related to data migrations or upgrades.
    - **`ProjectUpdate.md`**: Document detailing project updates.
    - **`Restructure.md`**: Document detailing project restructuring.
    - **`structure.md`**: Describes the project's directory structure (this file).
- **`internal/`**: Contains internal application logic, not intended for external consumption. Each sub-directory represents a distinct feature or service.
- **`pkg/`**: Contains shared libraries that are ok to be used by external applications.
    - **`core/`**: Core types for the application.
        - **`types.go`**: Defines core data structures.
- **`scripts/`**: Contains various utility scripts.
    - **`Build.sh`**: Script to build the application.
    - **`BuildAndRun.sh`**: Script to build and run the application.
    - **`BuildToWindows.sh`**: Script to build the application for Windows.
- **`vendor/`**: Go module dependencies. This directory is managed by `go modules` and contains vendored copies of external libraries.

## `internal/` Directory Breakdown

The `internal/` directory is central to the application's logic, organizing code by feature. This modular structure separates concerns, making the codebase easier to navigate and maintain. Each subdirectory typically contains a `client.go` for API interactions, `types.go` for data structures, and a `view.go` for rendering the UI component.

- **`clock/`**: Displays the current time.
    - **`view.go`**: Renders the clock component.
- **`common/`**: Shared utilities used across different internal packages.
    - **`time.go`**: Time-related utility functions.
- **`hosthatch/`**: Interacts with the HostHatch API.
- **`schedule/`**: Manages and displays schedule information.
- **`tailscale/`**: Interacts with the Tailscale API to get network device status.
- **`truenas/`**: Interacts with the TrueNAS API to get system status.
- **`upcloud/`**: Interacts with the UpCloud API.
- **`weather/`**: Fetches and displays weather information.

This structure aims to keep the codebase organized, making it easier to navigate, understand, and maintain.