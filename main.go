package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type Profile struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	SSHKey  string `json:"ssh_key"`
	Current bool   `json:"current"`
}

type Config struct {
	Profiles []Profile `json:"profiles"`
}

const (
	configDir  = ".config/gs"
	configFile = "profiles.json"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		switchProfile()
		return
	}

	switch args[0] {
	case "setup":
		setupFlow()
	case "list":
		listProfiles()
	case "edit":
		editProfile()
	case "rm", "remove":
		removeProfile()
	case "help", "-h", "--help":
		showHelp()
	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		showHelp()
		os.Exit(1)
	}
}

func getConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	configDirPath := filepath.Join(usr.HomeDir, configDir)
	configFilePath := filepath.Join(configDirPath, configFile)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDirPath, 0755); err != nil {
		return "", err
	}

	return configFilePath, nil
}

func loadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Profiles: []Profile{}}, nil
		}
		return nil, err
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func saveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, jsonData, 0644)
}

func switchProfile() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if len(config.Profiles) == 0 {
		fmt.Println("No profiles found. Run 'gs setup' to create your first profile.")
		return
	}

	if len(config.Profiles) == 1 {
		fmt.Println("Only one profile exists. Run 'gs setup' to create another profile.")
		return
	}

	// Find current profile and switch to next
	currentIndex := -1
	for i, profile := range config.Profiles {
		if profile.Current {
			currentIndex = i
			config.Profiles[i].Current = false
			break
		}
	}

	// If no current profile found, set first as current
	if currentIndex == -1 {
		currentIndex = 0
	} else {
		currentIndex = (currentIndex + 1) % len(config.Profiles)
	}

	config.Profiles[currentIndex].Current = true
	newProfile := config.Profiles[currentIndex]

	if err := updateGitConfig(newProfile); err != nil {
		fmt.Printf("Error updating git config: %v\n", err)
		os.Exit(1)
	}

	if err := updateSSHConfig(newProfile); err != nil {
		fmt.Printf("Error updating SSH config: %v\n", err)
		os.Exit(1)
	}

	if err := saveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	clearScreen()
	fmt.Printf("Switched to profile: \033[34m%s\033[0m (%s)\n", newProfile.Name, newProfile.Email)
}

func setupFlow() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Git Profile Setup ===")

	// Get name
	fmt.Print("Enter profile name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Get email
	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	// Get SSH key path
	fmt.Print("Enter path to SSH private key: ")
	sshKey, _ := reader.ReadString('\n')
	sshKey = strings.TrimSpace(sshKey)

	// Expand tilde to home directory
	if strings.HasPrefix(sshKey, "~") {
		usr, _ := user.Current()
		sshKey = strings.Replace(sshKey, "~", usr.HomeDir, 1)
	}

	// Validate SSH key exists
	if _, err := os.Stat(sshKey); os.IsNotExist(err) {
		fmt.Printf("SSH key not found at: %s\n", sshKey)
		os.Exit(1)
	}

	// Load existing config
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Check if profile with same name exists
	for _, profile := range config.Profiles {
		if profile.Name == name {
			fmt.Printf("Profile with name '%s' already exists.\n", name)
			os.Exit(1)
		}
	}

	// Create new profile
	newProfile := Profile{
		Name:    name,
		Email:   email,
		SSHKey:  sshKey,
		Current: len(config.Profiles) == 0, // First profile is current by default
	}

	// If this is not the first profile, make others non-current
	if len(config.Profiles) > 0 {
		for i := range config.Profiles {
			config.Profiles[i].Current = false
		}
		newProfile.Current = true
	}

	config.Profiles = append(config.Profiles, newProfile)

	// Save config
	if err := saveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	// Update git and SSH configs if this is the current profile
	if newProfile.Current {
		if err := updateGitConfig(newProfile); err != nil {
			fmt.Printf("Error updating git config: %v\n", err)
			os.Exit(1)
		}

		if err := updateSSHConfig(newProfile); err != nil {
			fmt.Printf("Error updating SSH config: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Profile '%s' created successfully!\n", name)
}

func listProfiles() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if len(config.Profiles) == 0 {
		fmt.Println("No profiles found. Run 'gs setup' to create your first profile.")
		return
	}

	fmt.Println("=== Git Profiles ===")
	for _, profile := range config.Profiles {
		current := ""
		if profile.Current {
			current = " \033[32m(current)\033[0m"
		}
		fmt.Printf("• %s <%s>%s\n", profile.Name, profile.Email, current)
		fmt.Printf("  SSH Key: %s\n", profile.SSHKey)
		fmt.Println()
	}
}

func editProfile() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if len(config.Profiles) == 0 {
		fmt.Println("No profiles found. Run 'gs setup' to create your first profile.")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	// Show profiles
	fmt.Println("=== Select Profile to Edit ===")
	for i, profile := range config.Profiles {
		fmt.Printf("%d. %s <%s>\n", i+1, profile.Name, profile.Email)
	}

	fmt.Print("Enter profile number: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var profileIndex int
	if _, err := fmt.Sscanf(input, "%d", &profileIndex); err != nil || profileIndex < 1 || profileIndex > len(config.Profiles) {
		fmt.Println("Invalid profile number.")
		os.Exit(1)
	}
	profileIndex-- // Convert to 0-based index

	profile := &config.Profiles[profileIndex]

	// Edit name
	fmt.Printf("Current name: %s\n", profile.Name)
	fmt.Print("New name (press Enter to keep current): ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name != "" {
		profile.Name = name
	}

	// Edit email
	fmt.Printf("Current email: %s\n", profile.Email)
	fmt.Print("New email (press Enter to keep current): ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)
	if email != "" {
		profile.Email = email
	}

	// Edit SSH key
	fmt.Printf("Current SSH key: %s\n", profile.SSHKey)
	fmt.Print("New SSH key path (press Enter to keep current): ")
	sshKey, _ := reader.ReadString('\n')
	sshKey = strings.TrimSpace(sshKey)
	if sshKey != "" {
		// Expand tilde
		if strings.HasPrefix(sshKey, "~") {
			usr, _ := user.Current()
			sshKey = strings.Replace(sshKey, "~", usr.HomeDir, 1)
		}

		// Validate SSH key exists
		if _, err := os.Stat(sshKey); os.IsNotExist(err) {
			fmt.Printf("SSH key not found at: %s\n", sshKey)
			os.Exit(1)
		}
		profile.SSHKey = sshKey
	}

	// Save config
	if err := saveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	// Update git and SSH configs if this is the current profile
	if profile.Current {
		if err := updateGitConfig(*profile); err != nil {
			fmt.Printf("Error updating git config: %v\n", err)
			os.Exit(1)
		}

		if err := updateSSHConfig(*profile); err != nil {
			fmt.Printf("Error updating SSH config: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Profile '%s' updated successfully!\n", profile.Name)
}

func removeProfile() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	if len(config.Profiles) == 0 {
		fmt.Println("No profiles found.")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	// Show profiles
	fmt.Println("=== Select Profile to Remove ===")
	for i, profile := range config.Profiles {
		current := ""
		if profile.Current {
			current = " (current)"
		}
		fmt.Printf("%d. %s <%s>%s\n", i+1, profile.Name, profile.Email, current)
	}

	fmt.Print("Enter profile number: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	var profileIndex int
	if _, err := fmt.Sscanf(input, "%d", &profileIndex); err != nil || profileIndex < 1 || profileIndex > len(config.Profiles) {
		fmt.Println("Invalid profile number.")
		os.Exit(1)
	}
	profileIndex-- // Convert to 0-based index

	profileToRemove := config.Profiles[profileIndex]

	// First confirmation
	fmt.Printf("Are you sure you want to remove profile '%s' <%s>? (y/N): ", profileToRemove.Name, profileToRemove.Email)
	confirm1, _ := reader.ReadString('\n')
	confirm1 = strings.ToLower(strings.TrimSpace(confirm1))

	if confirm1 != "y" && confirm1 != "yes" {
		fmt.Println("Aborted.")
		return
	}

	// Second confirmation
	fmt.Printf("This action cannot be undone. Are you absolutely sure? (y/N): ")
	confirm2, _ := reader.ReadString('\n')
	confirm2 = strings.ToLower(strings.TrimSpace(confirm2))

	if confirm2 != "y" && confirm2 != "yes" {
		fmt.Println("Aborted.")
		return
	}

	// Remove profile
	config.Profiles = append(config.Profiles[:profileIndex], config.Profiles[profileIndex+1:]...)

	// If removed profile was current, make first profile current
	if profileToRemove.Current && len(config.Profiles) > 0 {
		config.Profiles[0].Current = true
		if err := updateGitConfig(config.Profiles[0]); err != nil {
			fmt.Printf("Error updating git config: %v\n", err)
		}
		if err := updateSSHConfig(config.Profiles[0]); err != nil {
			fmt.Printf("Error updating SSH config: %v\n", err)
		}
	}

	// Save config
	if err := saveConfig(config); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Profile '%s' removed successfully!\n", profileToRemove.Name)
}

func updateGitConfig(profile Profile) error {
	// Set global git config
	if err := exec.Command("git", "config", "--global", "user.name", profile.Name).Run(); err != nil {
		return fmt.Errorf("failed to set git user.name: %v", err)
	}

	if err := exec.Command("git", "config", "--global", "user.email", profile.Email).Run(); err != nil {
		return fmt.Errorf("failed to set git user.email: %v", err)
	}

	return nil
}

func updateSSHConfig(profile Profile) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	configPath := filepath.Join(usr.HomeDir, ".ssh", "config")

	// Read existing SSH config
	content, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create new SSH config
			sshConfig := fmt.Sprintf(`Host github.com
    HostName github.com
    User git
    IdentityFile %s
`, profile.SSHKey)
			return os.WriteFile(configPath, []byte(sshConfig), 0600)
		}
		return err
	}

	lines := strings.Split(string(content), "\n")

	// Update IdentityFile for github.com
	inGithubSection := false
	updated := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "Host ") && strings.Contains(trimmed, "github.com") {
			inGithubSection = true
			continue
		}

		if inGithubSection && strings.HasPrefix(trimmed, "Host ") {
			inGithubSection = false
		}

		if inGithubSection && strings.Contains(line, "IdentityFile") {
			lines[i] = fmt.Sprintf("    IdentityFile %s", profile.SSHKey)
			updated = true
			break
		}
	}

	// If no github.com section found, append one
	if !updated {
		githubConfig := fmt.Sprintf(`
Host github.com
    HostName github.com
    User git
    IdentityFile %s`, profile.SSHKey)
		lines = append(lines, githubConfig)
	}

	updatedContent := strings.Join(lines, "\n")
	return os.WriteFile(configPath, []byte(updatedContent), 0600)
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func showHelp() {
	fmt.Println(`gs - Git Profile Switcher

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
    
    Profiles are stored in ~/.config/gs/profiles.json`)
}
