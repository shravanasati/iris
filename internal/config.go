package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var (
	SupportedResolutions = []string{"1024x768", "1600x900", "1920x1080", "3840x2160"}
)

type Configuration struct {
	RemoteSource            string   `json:"remote_source"`             // remote source to fetch wallpapers from
	SearchTerms             []string `json:"search_terms"`              // wallpaper search terms for unsplash
	Resolution              string   `json:"resolution"`                // wallpaper resolution, defaults to 1600x900
	ChangeWallpaper         bool     `json:"change_wallpaper"`          // whether change wallpaper after a duration
	ChangeWallpaperDuration string   `json:"change_wallpaper_duration"` // if wallpaper has to be changed, then after how many minutes
	WallpaperDirectory      string   `json:"wallpaper_directory"`       // use wallpapers from a user specified directory instead of unsplash
	SelectionType           string   `json:"selection_type"`            // directory wallpaper selection type, either sorted or random
	SaveWallpaper           bool     `json:"save_wallpaper"`            // whether to save the used wallpapers or not
	SaveWallpaperDirectory  string   `json:"save_wallpaper_directory"`  // directory to save the used wallpapers
	CheckForUpdates         bool     `json:"check_for_updates"`         // whether to check for updates
}

// todo add more remote sources - github, reddit

func (c *Configuration) WriteConfig() {
	configFilePath := filepath.Join(GetIrisDir(), "config.json")

	configFile, fer := os.Create(configFilePath)
	if fer != nil {
		fmt.Println("Unable to write config due to following error:", fer)
		os.Exit(1)
	}
	defer configFile.Close()

	if _, wer := configFile.Write(jsonify(c)); wer != nil {
		fmt.Println("Unable to write config due to following error:", wer)
		os.Exit(1)
	}
}

func (c *Configuration) Show() {
	searchTerms := strings.Join(c.SearchTerms, " ")

	data := [][]string{
		{"Search Terms", searchTerms},
		{"Remote Source", c.RemoteSource},
		{"Resolution", c.Resolution},
		{"Change Wallpaper", fmt.Sprintf("%v", c.ChangeWallpaper)},
		{"Change Wallpaper Duration", c.ChangeWallpaperDuration},
		{"Wallpaper Directory", c.WallpaperDirectory},
		{"Selection Type", c.SelectionType},
		{"Save Wallpaper", fmt.Sprintf("%v", c.SaveWallpaper)},
		{"Save Wallpaper Directory", c.SaveWallpaperDirectory},
		{"Check for Updates", fmt.Sprintf("%v", c.CheckForUpdates)},
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

	subDirs := []string{"wallpapers", "temp", "cache"}
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
		Resolution:              "1920x1080",
		ChangeWallpaper:         false,
		ChangeWallpaperDuration: "5m",
		WallpaperDirectory:      "",
		SelectionType:           "random",
		SaveWallpaper:           false,
		SaveWallpaperDirectory:  filepath.Join(GetIrisDir(), "wallpapers"),
		RemoteSource:            "unsplash",
		CheckForUpdates:         true,
	}

	return &defaultConfig
}

func ReadConfig() *Configuration {
	config := Configuration{}

	configFilePath := filepath.Join(GetIrisDir(), "config.json")

	if !CheckPathExists(configFilePath) {
		defaultConfig := getDefaultConfig()

		defaultConfig.WriteConfig()

		return defaultConfig
	}

	configContent := readFile(configFilePath)
	if e := json.Unmarshal([]byte(configContent), &config); e != nil {
		fmt.Println("Looks like the iris configuration is corrupted/broken, rewriting it with default values.")
		defaultConfig := getDefaultConfig()
		defaultConfig.WriteConfig()
		return defaultConfig
	}
	return &config
}

// todo rewrite existing config with new options
// in ReadConfig use prepareConfig function to unmarshal into a map and then add those
// options if they dont exist
