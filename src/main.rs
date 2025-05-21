use anyhow::{anyhow, Context, Result};
use clap::{Command};
use colored::Colorize;
use serde::{Deserialize, Serialize};
use std::fs::{self, File};
use std::io::{self, BufReader, Read, Write};
use std::path::{Path, PathBuf};
use std::process;

#[derive(Debug, Serialize, Deserialize, Clone)]
struct Profile {
    name: String,
    email: String,
    ssh_key: String,
    current: bool,
}

#[derive(Debug, Serialize, Deserialize)]
struct Config {
    profiles: Vec<Profile>,
}

const CONFIG_DIR: &str = ".config/gs";
const CONFIG_FILE: &str = "profiles.json";

fn main() -> Result<()> {
    let matches = Command::new("gs")
        .about("Switch between Git profiles instantly")
        .subcommand(Command::new("setup").about("Add new profile"))
        .subcommand(Command::new("list").about("Show all profiles"))
        .subcommand(Command::new("edit").about("Edit an existing profile"))
        .subcommand(
            Command::new("rm")
                .alias("remove")  // Set "remove" as an alias for "rm"
                .about("Remove a profile")
        )
        .get_matches();

    match matches.subcommand() {
        Some(("setup", _)) => setup_flow()?,
        Some(("list", _)) => list_profiles()?,
        Some(("edit", _)) => edit_profile()?,
        Some(("rm", _)) => remove_profile()?,  // Only need one match now
        None => switch_profile()?,
        _ => {
            // For any other command, show our custom help
            show_help();
            return Ok(());
        }
    }

    Ok(())
}

fn get_config_path() -> Result<PathBuf> {
    let home_dir = dirs::home_dir().ok_or_else(|| anyhow!("Could not find home directory"))?;
    let config_dir = home_dir.join(CONFIG_DIR);
    let config_path = config_dir.join(CONFIG_FILE);

    // Create config directory if it doesn't exist
    if !config_dir.exists() {
        fs::create_dir_all(&config_dir).context("Failed to create config directory")?;
    }

    Ok(config_path)
}

fn load_config() -> Result<Config> {
    let config_path = get_config_path()?;

    // If file doesn't exist, return empty config
    if !config_path.exists() {
        return Ok(Config { profiles: vec![] });
    }

    let file = File::open(config_path).context("Failed to open config file")?;
    let reader = BufReader::new(file);
    let config: Config = serde_json::from_reader(reader).context("Failed to parse config file")?;

    Ok(config)
}

fn save_config(config: &Config) -> Result<()> {
    let config_path = get_config_path()?;
    let json_data = serde_json::to_string_pretty(config).context("Failed to serialize config")?;
    fs::write(config_path, json_data).context("Failed to save config file")?;
    Ok(())
}

fn switch_profile() -> Result<()> {
    let mut config = load_config()?;

    if config.profiles.is_empty() {
        println!("No profiles found. Run 'gs setup' to create your first profile.");
        return Ok(());
    }

    if config.profiles.len() == 1 {
        println!("Only one profile exists. Run 'gs setup' to create another profile.");
        return Ok(());
    }

    // Find current profile and switch to next
    let mut current_index = 0;
    let mut found_current = false;

    for (i, profile) in config.profiles.iter_mut().enumerate() {
        if profile.current {
            profile.current = false;
            current_index = i;
            found_current = true;
            break;
        }
    }

    // If no current profile found, set first as current
    let new_index = if found_current {
        (current_index + 1) % config.profiles.len()
    } else {
        0
    };

    config.profiles[new_index].current = true;
    let new_profile = config.profiles[new_index].clone();

    update_git_config(&new_profile).context("Failed to update git config")?;
    update_ssh_config(&new_profile).context("Failed to update SSH config")?;
    save_config(&config)?;

    clear_screen();
    println!(
        "Switched to profile: {} ({})",
        new_profile.name.blue(),
        new_profile.email
    );

    Ok(())
}

fn setup_flow() -> Result<()> {
    println!("=== Git Profile Setup ===");

    // Get name
    print!("Enter profile name: ");
    io::stdout().flush()?;
    let mut name = String::new();
    io::stdin().read_line(&mut name)?;
    let name = name.trim().to_string();

    // Get email
    print!("Enter email: ");
    io::stdout().flush()?;
    let mut email = String::new();
    io::stdin().read_line(&mut email)?;
    let email = email.trim().to_string();

    // Get SSH key path
    print!("Enter path to SSH private key: ");
    io::stdout().flush()?;
    let mut ssh_key = String::new();
    io::stdin().read_line(&mut ssh_key)?;
    let mut ssh_key = ssh_key.trim().to_string();

    // Expand tilde to home directory
    if ssh_key.starts_with('~') {
        if let Some(home_dir) = dirs::home_dir() {
            ssh_key = ssh_key.replacen('~', home_dir.to_str().unwrap(), 1);
        }
    }

    // Validate SSH key exists
    if !Path::new(&ssh_key).exists() {
        return Err(anyhow!("SSH key not found at: {}", ssh_key));
    }

    // Load existing config
    let mut config = load_config()?;

    // Check if profile with same name exists
    for profile in &config.profiles {
        if profile.name == name {
            return Err(anyhow!("Profile with name '{}' already exists", name));
        }
    }

    // Create new profile
    let is_first_profile = config.profiles.is_empty();
    let new_profile = Profile {
        name,
        email,
        ssh_key,
        current: true, // New profile is set as current
    };

    // If this is not the first profile, make others non-current
    if !is_first_profile {
        for profile in &mut config.profiles {
            profile.current = false;
        }
    }

    // Update git and SSH configs for the new profile
    update_git_config(&new_profile)?;
    update_ssh_config(&new_profile)?;

    config.profiles.push(new_profile.clone());

    // Save config
    save_config(&config)?;

    println!("Profile '{}' created successfully!", new_profile.name);
    Ok(())
}

fn list_profiles() -> Result<()> {
    let config = load_config()?;

    if config.profiles.is_empty() {
        println!("No profiles found. Run 'gs setup' to create your first profile.");
        return Ok(());
    }

    println!("=== Git Profiles ===");
    for profile in &config.profiles {
        let current = if profile.current {
            " (current)".green().to_string()
        } else {
            String::new()
        };
        println!("• {} <{}>{}",
            profile.name,
            profile.email,
            current
        );
        println!("  SSH Key: {}", profile.ssh_key);
        println!();
    }

    Ok(())
}

fn edit_profile() -> Result<()> {
    // First, load the config and get necessary information
    let config = load_config()?;

    if config.profiles.is_empty() {
        println!("No profiles found. Run 'gs setup' to create your first profile.");
        return Ok(());
    }

    // Show profiles
    println!("=== Select Profile to Edit ===");
    for (i, profile) in config.profiles.iter().enumerate() {
        println!("{}. {} <{}>", i + 1, profile.name, profile.email);
    }

    // Get profile selection
    print!("Enter profile number: ");
    io::stdout().flush()?;
    let mut input = String::new();
    io::stdin().read_line(&mut input)?;
    let input = input.trim();

    let profile_index: usize = match input.parse::<usize>() {
        Ok(n) if n > 0 && n <= config.profiles.len() => n - 1,
        _ => return Err(anyhow!("Invalid profile number")),
    };

    // Clone the profile we want to edit and check if it's current
    let original_profile = config.profiles[profile_index].clone();
    let was_current = original_profile.current;
    let mut updated_profile = original_profile.clone();

    // Edit name
    println!("Current name: {}", updated_profile.name);
    print!("New name (press Enter to keep current): ");
    io::stdout().flush()?;
    let mut name = String::new();
    io::stdin().read_line(&mut name)?;
    let name = name.trim();
    if !name.is_empty() {
        updated_profile.name = name.to_string();
    }

    // Edit email
    println!("Current email: {}", updated_profile.email);
    print!("New email (press Enter to keep current): ");
    io::stdout().flush()?;
    let mut email = String::new();
    io::stdin().read_line(&mut email)?;
    let email = email.trim();
    if !email.is_empty() {
        updated_profile.email = email.to_string();
    }

    // Edit SSH key
    println!("Current SSH key: {}", updated_profile.ssh_key);
    print!("New SSH key path (press Enter to keep current): ");
    io::stdout().flush()?;
    let mut ssh_key = String::new();
    io::stdin().read_line(&mut ssh_key)?;
    let ssh_key = ssh_key.trim();
    
    if !ssh_key.is_empty() {
        let mut expanded_key = ssh_key.to_string();
        // Expand tilde
        if expanded_key.starts_with('~') {
            if let Some(home_dir) = dirs::home_dir() {
                expanded_key = expanded_key.replacen('~', home_dir.to_str().unwrap(), 1);
            }
        }

        // Validate SSH key exists
        if !Path::new(&expanded_key).exists() {
            return Err(anyhow!("SSH key not found at: {}", expanded_key));
        }
        
        updated_profile.ssh_key = expanded_key;
    }

    // Update git and SSH configs if this is the current profile
    if was_current {
        update_git_config(&updated_profile)?;
        update_ssh_config(&updated_profile)?;
    }

    // Now create a new config with the updated profile
    let mut new_config = config;
    new_config.profiles[profile_index] = updated_profile;

    // Save config
    save_config(&new_config)?;

    println!("Profile '{}' updated successfully!", new_config.profiles[profile_index].name);
    Ok(())
}

fn remove_profile() -> Result<()> {
    let mut config = load_config()?;

    if config.profiles.is_empty() {
        println!("No profiles found.");
        return Ok(());
    }

    // Show profiles
    println!("=== Select Profile to Remove ===");
    for (i, profile) in config.profiles.iter().enumerate() {
        let current = if profile.current { " (current)" } else { "" };
        println!("{}. {} <{}>{}",
            i + 1,
            profile.name,
            profile.email,
            current
        );
    }

    // Get profile selection
    print!("Enter profile number: ");
    io::stdout().flush()?;
    let mut input = String::new();
    io::stdin().read_line(&mut input)?;
    let input = input.trim();

    let profile_index: usize = match input.parse::<usize>() {
        Ok(n) if n > 0 && n <= config.profiles.len() => n - 1,
        _ => return Err(anyhow!("Invalid profile number")),
    };

    // Clone the profile information we need before removing it
    let profile_name = config.profiles[profile_index].name.clone();
    let was_current = config.profiles[profile_index].current;
    
    // First confirmation
    print!("Are you sure you want to remove profile '{}' <{}>? (y/N): ",
        profile_name,
        config.profiles[profile_index].email
    );
    io::stdout().flush()?;
    let mut confirm1 = String::new();
    io::stdin().read_line(&mut confirm1)?;
    let confirm1 = confirm1.trim().to_lowercase();

    if confirm1 != "y" && confirm1 != "yes" {
        println!("Aborted.");
        return Ok(());
    }

    // Second confirmation
    print!("This action cannot be undone. Are you absolutely sure? (y/N): ");
    io::stdout().flush()?;
    let mut confirm2 = String::new();
    io::stdin().read_line(&mut confirm2)?;
    let confirm2 = confirm2.trim().to_lowercase();

    if confirm2 != "y" && confirm2 != "yes" {
        println!("Aborted.");
        return Ok(());
    }

    // Remove profile
    config.profiles.remove(profile_index);

    // If removed profile was current, make first profile current
    if was_current && !config.profiles.is_empty() {
        config.profiles[0].current = true;
        update_git_config(&config.profiles[0])?;
        update_ssh_config(&config.profiles[0])?;
    }

    // Save config
    save_config(&config)?;

    println!("Profile '{}' removed successfully!", profile_name);
    Ok(())
}

fn update_git_config(profile: &Profile) -> Result<()> {
    // Set global git config
    process::Command::new("git")
        .args(["config", "--global", "user.name", &profile.name])
        .output()
        .context("Failed to set git user.name")?;

    process::Command::new("git")
        .args(["config", "--global", "user.email", &profile.email])
        .output()
        .context("Failed to set git user.email")?;

    Ok(())
}

fn update_ssh_config(profile: &Profile) -> Result<()> {
    let home_dir = dirs::home_dir().ok_or_else(|| anyhow!("Could not find home directory"))?;
    let ssh_dir = home_dir.join(".ssh");
    let config_path = ssh_dir.join("config");

    // Create .ssh directory if it doesn't exist
    if !ssh_dir.exists() {
        fs::create_dir_all(&ssh_dir).context("Failed to create .ssh directory")?;
    }

    // Read existing SSH config if it exists
    let content = if config_path.exists() {
        let mut file = File::open(&config_path).context("Failed to open SSH config")?;
        let mut content = String::new();
        file.read_to_string(&mut content)?;
        content
    } else {
        String::new()
    };

    let lines: Vec<&str> = content.lines().collect();
    let mut new_lines = Vec::new();
    let mut in_github_section = false;
    let mut updated = false;

    for line in &lines {
        let trimmed = line.trim();

        if trimmed.starts_with("Host ") && trimmed.contains("github.com") {
            in_github_section = true;
            new_lines.push(line.to_string());
            continue;
        }

        if in_github_section && trimmed.starts_with("Host ") {
            in_github_section = false;
        }

        if in_github_section && trimmed.contains("IdentityFile") {
            new_lines.push(format!("    IdentityFile {}", profile.ssh_key));
            updated = true;
        } else {
            new_lines.push(line.to_string());
        }
    }

    // If no github.com section found, append one
    if !updated {
        if !new_lines.is_empty() && !new_lines.last().unwrap().is_empty() {
            new_lines.push(String::new()); // Add empty line for spacing
        }
        
        new_lines.push("Host github.com".to_string());
        new_lines.push("    HostName github.com".to_string());
        new_lines.push("    User git".to_string());
        new_lines.push(format!("    IdentityFile {}", profile.ssh_key));
    }

    let updated_content = new_lines.join("\n");
    fs::write(&config_path, updated_content).context("Failed to write SSH config")?;
    
    // Set permissions
    #[cfg(unix)]
    {
        use std::os::unix::fs::PermissionsExt;
        let mut perms = fs::metadata(&config_path)?.permissions();
        perms.set_mode(0o600);
        fs::set_permissions(&config_path, perms)?;
    }

    Ok(())
}

fn clear_screen() {
    
    #[cfg(not(target_os = "windows"))]
    {
        let _ = process::Command::new("clear").status();
    }
}

fn show_help() {
    println!("gs - Git Profile Switcher

USAGE:
    gs               Switch between profiles (toggle if only 2 profiles)
    gs setup         Set up a new profile
    gs list          List all profiles
    gs edit          Edit an existing profile
    gs rm            Remove a profile
    gs help          Show this help message

DESCRIPTION:
    gs helps you manage multiple Git profiles for different accounts.
    Each profile includes a name, email, and SSH key.
    
    Profiles are stored in ~/.config/gs/profiles.json");
}
