# StatusApp

StatusApp is a terminal-based application written in Go, utilizing the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI (Terminal User Interface) framework. It provides a concise, real-time overview of various system statuses directly in your terminal.

## Features

-   **Digital Clock:** Displays the current time using stylized ASCII art (Figlet).
-   **Tailscale Status:** Shows the connection status and details of your Tailscale devices, including key expiry.
-   **Weather Information:** Presents current weather conditions for a configured location.
-   **Meeting Schedule:** Displays upcoming meetings from a local schedule file.

## Getting Started

### Prerequisites

To run StatusApp, you need:
-   Go (version 1.25.5 or later)
-   Access to Tailscale API (for device status)
-   Access to [WeatherAPI](https://www.weatherapi.com) (for weather information)

### Configuration

StatusApp uses environment variables for configuration, loaded from a `.env` file in the project root. Create a `.env` file with the following variables:

```
TAILSCALE_TAILNET_ID=your_tailnet_id
TAILSCALE_API_KEY=your_tailscale_api_key
TAILSCALE_API_KEY_ID=your_tailscale_api_key_id
WEATHERAPI_API_KEY=your_weatherapi_key
WEATHERAPI_LOCATION=your_location
WEATHER_ICON_PATH=file_path_weather_json
SCHEDULE_FILE_PATH=file_path
```
The `weather.json` file is currently in `assets/` but is necessary in runtime

## File Structure: Schedule

Delimiter: `##`
One meeting per line, on the format: 

```
start_time##Title##Room##end_time
```

**Example**

```
08:00##Stand-up##Microsoft Teams Meeting##08:15
```

**Note** 
There are some very specific parsing of the meeting room, in [loader.go](internal/schedule/loader.go) that will have to be updated for other users

### Installation and Running

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/manfredbjorlin/StatusApp.git
    cd StatusApp
    ```

2.  **Ensure dependencies are synced (if not already):**
    ```bash
    go mod tidy
    go mod vendor
    ```

3.  **Build and run the application:**
    ```bash
    go build -o statusapp ./cmd
    ./statusapp
    ```

## Documentation

For a more in-depth understanding of the project:

-   [Codebase Overview](docs/codebase.md)
-   [Project Structure](docs/structure.md)
-   [Migration History](docs/migration.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
