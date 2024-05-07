package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

type Profile struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Current bool   `json:"current"`
}

type ProfilesList struct {
	Profiles []Profile
}

func main() {
	// Open profiles for reading
	file, err := os.Open("/Users/ianeblack/.profiles.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read profiles from file
	jsonData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading JSON data from file:", err)
		return
	}

	var myProfiles ProfilesList
	err = json.Unmarshal(jsonData, &myProfiles)
	if err != nil {
		fmt.Println("Error parsing JSON data:", err)
		return
	}

  // Switched to
  var currentProfile Profile

	// Hot-swap profile depending on current
	for i := range myProfiles.Profiles {
		myProfiles.Profiles[i].Current = !myProfiles.Profiles[i].Current
    if myProfiles.Profiles[i].Current == true {
      updateConfig(myProfiles.Profiles[i])
      currentProfile = myProfiles.Profiles[i]
    }
	}

	// Marshal the updated data to JSON
	updatedJsonData, err := json.MarshalIndent(myProfiles, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling updated JSON data:", err)
		return
	}

	// Open file for writing (this will truncate the file if it exists)
	file, err = os.OpenFile("/Users/ianeblack/.profiles.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file for writing:", err)
		return
	}
	defer file.Close()

	// Write updated JSON data back to file
	_, err = file.Write(updatedJsonData)
	if err != nil {
		fmt.Println("Error writing updated JSON data to file:", err)
		return
	}

  fmt.Println("\nHello \033[34m" + currentProfile.Name + "\033[0m\n")

}

func updateConfig(profile Profile) {
	// Open ssh config
	config_file, err := os.Open("/Users/ianeblack/.ssh/config")
	if err != nil {
		fmt.Println("Error opening config:", err)
	}
	defer config_file.Close()

	// Read ssh config_file
	configData, err := io.ReadAll(config_file)
	if err != nil {
		fmt.Println("error reading config:", err)
	}

	configLines := strings.Split(string(configData), "\n")

	// Update the IdentityFile line in the SSH config
	for i, line := range configLines {
		if strings.Contains(line, "IdentityFile") {
			configLines[i] = fmt.Sprintf("IdentityFile %s", profile.Path)
			break
		}
	}
	updatedConfig := strings.Join(configLines, "\n")

	// Write the updated SSH config back to the file
	if err := os.WriteFile("/Users/ianeblack/.ssh/config", []byte(updatedConfig), fs.FileMode(0644)); err != nil {
		fmt.Println("Error writing updated SSH config:", err)
		return
	}

}
