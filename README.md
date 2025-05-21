# gs

Switch between Git profiles instantly.

## Features

- Easily manage multiple Git profiles for different accounts
- Each profile includes name, email, and SSH key
- Simple command-line interface
- Automatic Git and SSH configuration updates

## Installation

### From Source

```bash
git clone https://github.com/yourusername/git-switch-rs gs
cd gs
cargo build --release
sudo cp target/release/gs /usr/local/bin/
```

### Using Cargo

```bash
cargo install --git https://github.com/yourusername/git-switch-rs
```

## Usage

```bash
gs          # Switch between profiles
gs setup    # Add new profile
gs list     # Show all profiles
gs edit     # Edit an existing profile
gs rm       # Remove a profile
```

## How It Works

Each profile contains:
- Name: Your Git username
- Email: Your Git email address
- SSH key: Path to your SSH private key

When you switch profiles, `gs` updates:
1. Your global Git configuration
2. Your SSH configuration for GitHub

## Configuration

Profiles are stored in `~/.config/gs/profiles.json`.

## License

MIT
