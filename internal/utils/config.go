package utils

import (
    "os"
    "strings"
    "log"

    "gopkg.in/ini.v1"
)

// GetProfiles parses the ~/.databrickscfg file and returns the list of profiles.
func GetProfiles(configPath string) ([]string) {
    homeDir, _ := os.UserHomeDir()
    configPathFixed := strings.Replace(configPath, "~", homeDir, 1)

    // Load the config file.
    cfg, err := ini.Load(configPathFixed)
    if err != nil {
        log.Fatalf("Failed to fetch profiles: %v", err)
    }

    // Get the list of profiles (sections in the ini file).
    profiles := cfg.SectionStrings()

    return profiles
}
