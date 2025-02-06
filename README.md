# Sanity

Solana vanity address generator with a beautiful and intuitive CLI & TUI. Generate custom Solana wallet addresses with your desired prefix.

## Features

- Intuitive Terminal User Interface (TUI)
- Command Line Interface (CLI) support
- Concurrent address generation
- Customizable search parameters
- Automatic private key saving
- Timeout configuration

## Installation

```bash
go install github.com/tmlnv/sanity@latest
```

Or clone and build from source:

```bash
git clone https://github.com/tmlnv/sanity.git
cd sanity
go build -o sanity ./cmd/sanity
```

## Usage

### TUI Mode

Simply run the program without any flags to enter TUI mode:

```bash
sanity
```

In TUI mode, you can:

1. Enter your desired address prefix
2. Specify the number of addresses to generate
3. Set the number of concurrent threads (defaults to CPU cores if left empty)
4. Configure a timeout duration (e.g., "30s", "5m", "1h")

Use Tab/Shift+Tab or Up/Down arrows to navigate between fields.

### CLI Mode

Run with flags for CLI mode:

```bash
sanity -prefix <desired_prefix> [-n <number_of_addresses>] [-t <threads>] [-timeout <duration>]
```

Example:

```bash
sanity -prefix abc -n 1 -t 4 -timeout 5m
```

## Security

Generated private keys are automatically saved to a file in your current directory. Keep these keys secure and never share them.

## Disclaimer

This tool is for educational and experimental purposes only. Please be aware that:

- Generated addresses and private keys should be thoroughly verified before use
- Always follow best security practices when handling cryptocurrency wallets
- The author is not responsible for anything

## License

MIT License - see LICENSE file for details

## TODO

- [x] Validate address input to be compatible with Solana to not
      waste resources trying to find the impossible one.
- [x] Suffix
- [x] Regexp
- [x] Timeout input proper validation.
- [ ] Validate timeout from pure number for CLI
- [ ] Case sensitivity
- [ ] CLI help
