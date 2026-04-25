package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Configuration struct {
	RemoteSource            string   `json:"remote_source"`             // remote source to fetch wallpapers from
	SearchTerms             []string `json:"search_terms"`              // wallpaper search terms
	ChangeWallpaper         bool     `json:"change_wallpaper"`          // whether change wallpaper after a duration
	ChangeWallpaperDuration string   `json:"change_wallpaper_duration"` // if wallpaper has to be changed, then after how many minutes
	WallpaperFile           string   `json:"wallpaper_file"`            // path to the wallpaper file, might be a video as well
	WallpaperDirectory      string   `json:"wallpaper_directory"`       // use wallpapers from a user specified directory
	SelectionType           string   `json:"selection_type"`            // directory wallpaper selection type, either sorted or random
	SaveWallpaper           bool     `json:"save_wallpaper"`            // whether to save the used wallpapers or not
	SaveWallpaperDirectory  string   `json:"save_wallpaper_directory"`  // directory to save the used wallpapers
	CheckForUpdates         bool     `json:"check_for_updates"`         // whether to check for updates
	GitHubAPIToken          string   `json:"github_api_token"`          // github api token to perform auth requests
}

func (c *Configuration) WriteConfig() {
	configFilePath := filepath.Join(GetIrisDir(), "config.json")
	LogInfof("config", "writing configuration to: %s", configFilePath)

	configFile, fer := os.Create(configFilePath)
	if fer != nil {
		LogErrorf("config", "failed to create config file: %v", fer)
		fmt.Println("Unable to write config due to following error:", fer)
		os.Exit(1)
	}
	defer configFile.Close()

	if _, wer := configFile.Write(jsonify(c)); wer != nil {
		LogErrorf("config", "failed to write config data: %v", wer)
		fmt.Println("Unable to write config due to following error:", wer)
		os.Exit(1)
	}
	LogInfof("config", "configuration written successfully")
}

func (c *Configuration) Show() {
	searchTerms := strings.Join(c.SearchTerms, " ")

	data := [][]string{
		{"Search Terms", searchTerms},
		{"Remote Source", c.RemoteSource},
		{"Change Wallpaper", fmt.Sprintf("%v", c.ChangeWallpaper)},
		{"Change Wallpaper Duration", c.ChangeWallpaperDuration},
		{"Wallpaper File", c.WallpaperFile},
		{"Wallpaper Directory", c.WallpaperDirectory},
		{"Selection Type", c.SelectionType},
		{"Save Wallpaper", fmt.Sprintf("%v", c.SaveWallpaper)},
		{"Save Wallpaper Directory", c.SaveWallpaperDirectory},
		{"Check for Updates", fmt.Sprintf("%v", c.CheckForUpdates)},
		{"GitHub API Token", c.GitHubAPIToken},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Option", "Value"})
	table.AppendBulk(data)

	table.Render()
}

// GetIrisDir returns the iris home directory, namely `~/.iris`. Also creates the directory if it doesn't exists, and the necessary subfolders wallpapers, temp and cache.
func GetIrisDir() string {
	usr, e := user.Current()
	if e != nil {
		panic(e)
	}

	// * determining iris's directory
	dir := filepath.Join(usr.HomeDir, ".iris")

	if !CheckPathExists(dir) {
		os.Mkdir(dir, os.ModePerm)
	}

	subDirs := []string{"wallpapers", "temp", "cache", "logs"}
	for _, subDir := range subDirs {
		dirPath := filepath.Join(dir, subDir)
		if !CheckPathExists(dirPath) {
			os.Mkdir(dirPath, os.ModePerm)
		}
	}

	return dir
}

func getDefaultConfig() *Configuration {
	defaultConfig := Configuration{
		SearchTerms:             []string{"nature"},
		ChangeWallpaper:         false,
		ChangeWallpaperDuration: "5m",
		WallpaperFile:           "",
		WallpaperDirectory:      "",
		SelectionType:           "random",
		SaveWallpaper:           false,
		SaveWallpaperDirectory:  filepath.Join(GetIrisDir(), "wallpapers"),
		RemoteSource:            "spotlight",
		CheckForUpdates:         true,
		GitHubAPIToken:          "",
	}

	return &defaultConfig
}

func ReadConfig() *Configuration {
	config := Configuration{}

	configFilePath := filepath.Join(GetIrisDir(), "config.json")
	LogInfof("config", "reading configuration from: %s", configFilePath)

	if !CheckPathExists(configFilePath) {
		LogInfof("config", "config file not found, creating default")
		defaultConfig := getDefaultConfig()

		defaultConfig.WriteConfig()

		return defaultConfig
	}

	configContent := readFile(configFilePath)

	// implement backward compatibility
	// unmarshal into a map to see which keys are present
	var configMap map[string]any
	if err := json.Unmarshal([]byte(configContent), &configMap); err != nil {
		LogErrorf("config", "failed to parse config json: %v", err)
		fmt.Printf("unable to read config: %v\n", err)
		fmt.Println("Looks like the iris configuration is corrupted/broken, rewriting it with default values.")
		defaultConfig := getDefaultConfig()
		defaultConfig.WriteConfig()
		return defaultConfig
	}

	defaultConfig := getDefaultConfig()
	serializedDefault, _ := json.Marshal(defaultConfig)
	var defaultMap map[string]any
	json.Unmarshal(serializedDefault, &defaultMap)

	needsUpdate := false
	for key, value := range defaultMap {
		actualValue, exists := configMap[key]
		if !exists {
			LogInfof("config", "adding missing key: %s", key)
			configMap[key] = value
			needsUpdate = true
		} else {
			// type check to ensure existing values match the expected type
			if reflect.TypeOf(actualValue) != reflect.TypeOf(value) {
				LogWarnf("config", "type mismatch for key %s, resetting to default", key)
				configMap[key] = value
				needsUpdate = true
			}
		}
	}

	// remove keys that are no longer supported (like resolution)
	for key := range configMap {
		if _, exists := defaultMap[key]; !exists {
			LogInfof("config", "removing unsupported key: %s", key)
			delete(configMap, key)
			needsUpdate = true
		}
	}

	if needsUpdate {
		LogInfof("config", "updating config file with new keys/types")
		// unmarshal the updated map back into the config struct
		updatedConfigBytes, _ := json.Marshal(configMap)
		json.Unmarshal(updatedConfigBytes, &config)
		config.WriteConfig()
	} else {
		// if no update was needed, just unmarshal the original content into the struct
		if e := json.Unmarshal([]byte(configContent), &config); e != nil {
			LogErrorf("config", "failed to unmarshal final config: %v", e)
			fmt.Printf("unable to read config: %v\n", e)
			defaultConfig := getDefaultConfig()
			defaultConfig.WriteConfig()
			return defaultConfig
		}
		LogInfof("config", "configuration loaded successfully")
	}

	return &config
}
