# Stamp

[![Release](https://img.shields.io/github/v/release/totocaster/stamp)](https://github.com/totocaster/stamp/releases)
[![CI](https://github.com/totocaster/stamp/actions/workflows/ci.yml/badge.svg)](https://github.com/totocaster/stamp/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/totocaster/stamp)](https://go.dev/)
[![License](https://img.shields.io/github/license/totocaster/stamp)](https://github.com/totocaster/stamp/blob/main/LICENSE)

A simple Go CLI tool for generating note filenames based on date/time following Toto's note naming conventions. Available as both `stamp` and `nid` commands.

## Features

- 📅 **Multiple Note Types**: Daily, fleeting, voice, analog/slipbox, monthly, yearly, and project notes
- 🔢 **Smart Counters**: Automatic sequential numbering for analog (daily reset) and project notes (persistent)
- ⚙️ **Configurable**: YAML configuration for timezone, defaults, and counter storage
- 📋 **Clipboard Support**: Copy generated names directly to clipboard (macOS)
- 🚀 **Fast & Lightweight**: Written in Go for instant execution
- 🔄 **Dual Commands**: Use as `stamp` or `nid` (Note ID)
- 🧭 **Obsidian-Aware**: Automatically picks up [Daily Notes](https://help.obsidian.md/Plugins/Core+plugins/Daily+notes) and [Unique Note Creator](https://github.com/adriano-tirloni/unique-note-creator) formats when run inside a vault

## Quick Start

```bash
# Default timestamp (YYYY-MM-DD-HHMM)
$ stamp
2025-11-12-1534

# Daily note
$ stamp daily
2025-11-12

# Fleeting note with timestamp
$ stamp fleeting
2025-11-12-F153045

# Project with auto-increment
$ stamp project
P0395

$ stamp project "New CLI Tool"
P0396 New CLI Tool

# Analog note (sequential per day)
$ stamp analog
2025-11-12-A1
$ stamp analog
2025-11-12-A2
```

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap totocaster/tap
brew install stamp  # Installs both stamp and nid commands
```

### Direct Download

Download the latest binary for your platform from the [releases page](https://github.com/totocaster/stamp/releases).

### Using Go

```bash
go install github.com/toto/stamp/cmd/stamp@latest
```

### From Source

```bash
# Clone the repository
git clone https://github.com/totocaster/stamp.git
cd stamp

# Build and install using Make
make install

# Or build manually
go build -o stamp cmd/stamp/main.go
sudo cp stamp /usr/local/bin/
sudo ln -s /usr/local/bin/stamp /usr/local/bin/nid
```

## Usage

### Note Types

| Type | Format | Example | Description |
|------|--------|---------|-------------|
| Default | `YYYY-MM-DD-HHMM` | `2025-11-12-1534` | Default timestamp (24-hour format) |
| Daily | `YYYY-MM-DD` | `2025-11-12` | Daily notes |
| Fleeting | `YYYY-MM-DD-FHHMMSS` | `2025-11-12-F153045` | Quick capture with seconds |
| Voice | `YYYY-MM-DD-VTHHMMSS` | `2025-11-12-VT153045` | Voice transcripts with seconds |
| Analog | `YYYY-MM-DD-AN` | `2025-11-12-A3` | Sequential slipbox notes (daily reset) |
| Monthly | `YYYY-MM` | `2025-11` | Monthly reviews |
| Yearly | `YYYY` | `2025` | Yearly reviews |
| Project | `PXXXX [title]` | `P0395 New Project` | Project numbers (persistent) |

### Flags

```bash
# Add .md extension
$ stamp --ext
2025-11-12-1534.md

# Copy to clipboard (macOS)
$ stamp --copy
2025-11-12-1534
Copied to clipboard!

# Quiet mode (no extra output)
$ stamp -q --copy
2025-11-12-1534

# Combine multiple flags
$ stamp daily --ext --copy
2025-11-12.md
Copied to clipboard!
```

### Counter Management

For sequential types (analog, project), manage counters with:

```bash
# Check next number without incrementing
$ stamp analog --check
2025-11-12-A3

$ stamp project --check
P0397

# Reset counter
$ stamp analog --reset
Counter reset for analog notes

# Set specific project number
$ stamp project --set 500
Project counter set to 500

# Show current counter value
$ stamp analog --counter
Current analog counter for 2025-11-12: 2

$ stamp project --counter
Current project counter: 397
```

## Configuration

Optional configuration file at `~/.stamp/config.yaml`:

```yaml
# Timezone for timestamps (default: system timezone)
timezone: "Asia/Tokyo"

# Always add .md extension
always_extension: false

# Counter storage location
counter_file: "~/.stamp/counters.json"

# Starting project number
project_start: 395
```

### Obsidian Integration

When `stamp` runs inside an Obsidian vault it mirrors your existing date formats.

- **Vault detection**: the CLI walks up from the current working directory until it finds a `.obsidian/` folder.
- **Daily Notes**: if the core plugin is enabled in `.obsidian/core-plugins.json`, `stamp` reads `daily-notes.json` (or `dailyNotes.format` within `app.json`) and translates the Moment-style string to Go's layout before emitting daily filenames.
- **Unique Note Creator**: when the community plugin is enabled (or its folder exists) the tool inspects `.obsidian/plugins/unique-note-creator/data.json` for filename patterns and uses them for the default command.
- **Graceful fallback**: missing files or unsupported tokens leave `stamp` on its built-in formats, and any read/parse issues are emitted as warnings on stderr without interrupting execution.

## Examples

### Daily Workflow

```bash
# Morning daily note
$ stamp daily --ext --copy
2025-11-12.md
Copied to clipboard!

# Quick fleeting thought
$ stamp fleeting
2025-11-12-F093045

# New analog note in sequence
$ stamp analog
2025-11-12-A1

# Start new project
$ stamp project "Stamp CLI Tool"
P0395 Stamp CLI Tool
```

### Counter Workflow

```bash
# Check what's next without using it
$ stamp project --check
P0397

# Actually create the project
$ stamp project "New Feature"
P0397 New Feature

# Made a mistake? Reset
$ stamp project --set 396
Project counter set to 396
```

## Development

### Project Structure

```
stamp/
├── cmd/stamp/          # Main application entry
├── internal/           # Internal packages
│   ├── config/         # Configuration handling
│   ├── counter/        # Counter management
│   ├── generator/      # Timestamp generation
│   └── clipboard/      # Clipboard operations
├── Makefile            # Build automation
├── README.md           # Documentation
├── LICENSE             # MIT License
└── go.mod              # Go module definition
```

### Building

```bash
# Build binary
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Build for multiple platforms
make release-build

# Format code
make fmt

# Run linter
make lint
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/generator
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) file for details

## Author

Created by Tornike (Toto) Tvalavadze for personal note-taking workflows.
