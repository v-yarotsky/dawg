package dawg

import (
	"fmt"
	"os"
	"path/filepath"
)

func SettingsDirPath() (string, error) {
	preferencesRoot := os.Getenv("alfred_preferences")
	if preferencesRoot == "" {
		return "", fmt.Errorf("aflred_preferences is not set. Are you running this manually?")
	}
	workflowsPath := filepath.Join(preferencesRoot, "workflows")
	settingsDir := filepath.Join(workflowsPath, fmt.Sprintf("%s.%s", BundleID, "syncedsettings"))
	return settingsDir, nil
}

func SettingsFilePath() (string, error) {
	if settingsDir, err := SettingsDirPath(); err != nil {
		return "", err
	} else {
		return filepath.Join(settingsDir, "dawg.json"), nil
	}
}
