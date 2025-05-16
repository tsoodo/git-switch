# GS - Git Profile Switcher

A command-line tool for managing multiple Git profiles with different SSH keys, perfect for developers who work with multiple GitHub accounts.

## Features

- **Quick Profile Switching**: Toggle between profiles with a single command
- **Multiple Account Support**: Manage unlimited Git profiles
- **Automatic Configuration**: Updates both Git config and SSH config automatically
- **Easy Setup**: Interactive setup flow for new profiles
- **Profile Management**: List, edit, and remove profiles with simple commands

## Installation

### Homebrew (Recommended)

```bash
# Install directly from the repository
brew install tsoodo/git-switch/gs
```

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/tsoodo/git-switch/releases)
2. Extract the binary to your PATH (e.g., `/usr/local/bin/`)
3. Make it executable: `chmod +x gs`

## Configuration File Location

- **macOS/Linux**: `~/.config/gs/profiles.json`

## Requirements

- Git installed on your system
- SSH keys generated for your GitHub accounts
- SSH config file at `~/.ssh/config` (created automatically if needed)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details
