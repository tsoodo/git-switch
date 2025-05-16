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
# Add the tap (replace with your actual repo)
brew tap YOUR_USERNAME/gs

# Install gs
brew install gs
```

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/YOUR_USERNAME/gs/releases)
2. Extract the binary to your PATH (e.g., `/usr/local/bin/`)
3. Make it executable: `chmod +x gs`

## Usage

### Basic Commands

```bash
# Switch between profiles (toggles if only 2 profiles)
gs

# Set up a new profile (interactive)
gs setup

# List all profiles
gs list

# Edit an existing profile
gs edit

# Remove a profile (with confirmation)
gs rm

# Show help
gs help
```

### Setup Flow

When you run `gs setup`, you'll be prompted for:
- **Profile Name**: A friendly name for the profile (e.g., "Personal", "Work")
- **Email**: The email associated with this Git profile
- **SSH Key Path**: Path to the private SSH key for this profile

### Configuration

Profiles are stored in `~/.config/gs/profiles.json`. Each profile includes:
- Name and email for Git commits
- SSH private key path for GitHub authentication

## How It Works

When you switch profiles, `gs` automatically:
1. Updates `git config --global user.name` and `user.email`
2. Updates `~/.ssh/config` to use the correct `IdentityFile` for github.com
3. Marks the new profile as current in the configuration

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
