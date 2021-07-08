package main

import (
	"bufio"
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

var (
	supportedResolutions = []string{"1024x768", "1600x900", "1920x1080", "3840x2160"}
)

type Configuration struct {
	SearchTerms             []string `json:"search_terms"`              // wallpaper search terms for unsplash
	Resolution              string   `json:"resolution"`                // wallpaper resolution, defaults to 1600x900
	ChangeWallpaper         bool     `json:"change_wallpaper"`          // whether change wallpaper after a duration
	ChangeWallpaperDuration int      `json:"change_wallpaper_duration"` // if wallpaper has to be changed, then after how many minutes
	WallpaperDirectory      string   `json:"wallpaper_directory"`       // use wallpapers from a user specified directory instead of unsplash
	SelectionType           string   `json:"selection_type"`            // directory wallpaper selection type, either sorted or random
	SaveWallpaper           bool     `json:"save_wallpaper"`            // whether to save the used wallpapers or not
}

// readFile reads the given file and returns the string content of the same.
func readFile(file string) string {
	f, ferr := os.Open(file)
	if ferr != nil {
		panic(ferr)
	}
	defer f.Close()

	text := ""
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text += scanner.Text()
	}

	return text
}

// getIrisDir returns the iris home directory, namely `~/.iris`. Also creates the directory if it doesnt exists.
func getIrisDir() string {
	usr, e := user.Current()
	if e != nil {
		panic(e)
	}

	// * determining iris's directory
	dir := filepath.Join(usr.HomeDir, ".iris")

	_, er := os.Stat(dir)
	if os.IsNotExist(er) {
		os.Mkdir(dir, os.ModePerm)
	}

	wallpaperDir := filepath.Join(dir, "wallpapers")
	_, err := os.Stat(wallpaperDir)
	if os.IsNotExist(err) {
		os.Mkdir(wallpaperDir, os.ModePerm)
	}

	return dir
}

func checkFileExists(filePath string) bool {
	_, e := os.Stat(filePath)
	return !os.IsNotExist(e)
}

func jsonifyConfig(config *Configuration) []byte {
	byteArray, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		panic(err)
	}
	return (byteArray)
}

func readConfig() *Configuration {
	config := Configuration{}

	configFilePath := filepath.Join(getIrisDir(), "config.json")

	if !checkFileExists(configFilePath) {
		defaultConfig := Configuration{
			SearchTerms:             []string{"nature"},
			Resolution:              "1600x900",
			ChangeWallpaper:         false,
			ChangeWallpaperDuration: -1,
			WallpaperDirectory:      "",
			SelectionType:           "random",
			SaveWallpaper:           false,
		}

		configFile, fer := os.Create(configFilePath)
		if fer != nil {
			panic(fer)
		}
		defer configFile.Close()
		if _, wer := configFile.Write(jsonifyConfig(&defaultConfig)); wer != nil {
			panic(wer)
		}

		return &defaultConfig
	}

	configContent := readFile(configFilePath)
	if e := json.Unmarshal([]byte(configContent), &config); e != nil {
		panic(e)
	}

	return &config
}
