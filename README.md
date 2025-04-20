# Sanity

[![Build](https://github.com/tmlnv/sanity/actions/workflows/build.yml/badge.svg)](https://github.com/tmlnv/sanity/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/tmlnv/sanity/graph/badge.svg?token=PR3E1MGGRZ)](https://codecov.io/gh/tmlnv/sanity)

Solana vanity address generator with CLI & TUI.
Generate custom Solana wallet addresses with your desired prefix, suffix or regexp.

![demo_tui](https://tmlnv.com/sanity_demo.gif)

## Features

- Intuitive Terminal User Interface (TUI)
- Command Line Interface (CLI) support
- Concurrent address generation
- Customizable search parameters
- Automatic private key saving
- Timeout configuration

## Installation

You can install and run Sanity either by using Docker or by cloning the repository and building it directly.

### Using Docker

Pull the image from Docker Hub:

```bash
docker pull tmlnv/sanity:latest
```

Run the application using Docker:

```bash
# To save generated keys to your current host directory, mounting a volume is necessary.
# The default PrivateKeysFile is 'sanity.private.log' in the container's working directory.

# Run in TUI mode (interactive)
# Note: TUI rendering (colors, cursor) might vary depending on your terminal and Docker setup and not provide a full experience using Docker.
docker run -it --rm -v "$(pwd):/app" tmlnv/sanity:latest

# Run in CLI mode with flags
docker run --rm -v "$(pwd):/app" tmlnv/sanity:latest -prefix 333 -count 1
```

### From Source

To build and run from source, you need Go installed.

1. Clone the repository:

```bash
git clone https://github.com/tmlnv/sanity.git
cd sanity
```

2. Build the executable:

```bash
go build -o sanity .
```

3. Run the application:

```bash
# Run in TUI mode
./sanity

# Run in CLI mode with flags
./sanity -prefix 333 -count 1
```

## Usage

Sanity can be used in two modes: TUI (Terminal User Interface) or CLI (Command Line Interface).

### TUI Mode

Simply run the program without any flags to enter TUI mode:

```bash
./sanity
```

In TUI mode, you can:

1. Enter your desired address prefix
2. Specify the number of addresses to generate
3. Set the number of concurrent workers (defaults to CPU cores if left empty)
4. Configure a timeout duration (e.g., "30s", "5m", "1h")

Use Tab/Shift+Tab or Up/Down arrows to navigate between fields.

### CLI Mode

Run with flags for CLI mode:

```bash
./sanity [-prefix <prefix>] [-suffix <suffix>] [-regexp <pattern>] [-count <number>] [-workers <workers>] [-timeout <duration>]
```

Examples:

```bash
# Generate address with prefix
sanity -prefix 123 -count 1 -workers 4 -timeout 5m

# Generate address with suffix
sanity -suffix 123 -count 2

# Generate address with both prefix and suffix
sanity -prefix 123 -suffix 321

# Generate address matching regular expression
sanity -regexp '^123.*321$' -timeout 10m
```

Options:

- `-prefix`: Desired prefix for the address
- `-suffix`: Desired suffix for the address
- `-regexp`: Regular expression pattern to match
- `-count`: Number of addresses to generate (default: 1)
- `-workers`: Number of concurrent workers (default: CPU cores)
- `-timeout`: Maximum duration to search (e.g., "30s", "5m", "1h")

## Security

Generated private keys are automatically saved to a file in your current directory.
Keep these keys secure and never share them.

## Disclaimer

This tool is for educational and experimental purposes only. Please be aware that:

- Generated addresses and private keys should be thoroughly verified before use
- Always follow best security practices when handling cryptocurrency wallets
- The author is not responsible for anything

## TODO

- [x] Validate address input to be compatible with Solana to not
      waste resources trying to find the impossible one.
- [x] Suffix
- [x] Regexp
- [x] Timeout input proper validation.
- [x] Validate timeout from pure number for CLI
- [x] Case sensitivity
- [x] CLI help
- [x] Tests
