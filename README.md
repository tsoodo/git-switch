[![Rust](https://github.com/tsoodo/git-switch/actions/workflows/rust.yml/badge.svg)](https://github.com/tsoodo/git-switch/actions/workflows/rust.yml)


# gs

Switch between Git profiles instantly. A command-line tool that helps
developers who work with multiple Git accounts (personal, work, open source)
manage their identities seamlessly.

## Why use gs?

Managing multiple Git identities can be frustrating. If you've ever
accidentally committed to a work project with your personal email or vice
versa, you know the pain. Git's global config is convenient but problematic
when switching contexts.

`gs` solves this by letting you:

- Define multiple Git profiles with different names, emails, and SSH keys
- Switch between them with a single command
- Automatically update both Git and SSH configurations
- Keep your commit history clean and correctly attributed

Perfect for developers who:
- Work on both company and personal projects
- Contribute to open source under different identities
- Manage multiple clients or organizations
- Need to maintain separate SSH keys for different services

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
