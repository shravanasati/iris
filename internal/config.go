package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
)

var (
	SupportedResolutions = []string{"1024x768", "1600x900", "1920x1080", "3840x2160"}
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

func (c *Configuration) WriteConfig() {
	configFilePath := filepath.Join(getIrisDir(), "config.json")

	configFile, fer := os.Create(configFilePath)
	if fer != nil {
		fmt.Println("Unable to write config due to following error:", fer)
		os.Exit(1)
	}
	defer configFile.Close()

	if _, wer := configFile.Write(jsonifyConfig(c)); wer != nil {
		fmt.Println("Unable to write config due to following error:", wer)
		os.Exit(1)
	}
}

func (c *Configuration) Show() {
	searchTerms := ""
	for _, term := range c.SearchTerms {
		searchTerms += term + " "
	}

	data := [][]string{
		{"Search Terms", searchTerms},
		{"Resolution", c.Resolution},
		{"Change Wallpaper", fmt.Sprintf("%v", c.ChangeWallpaper)},
		{"Change Wallpaper Duration", fmt.Sprintf("%v minute(s)", c.ChangeWallpaperDuration)},
		{"Wallpaper Directory", c.WallpaperDirectory},
		{"Selection Type", c.SelectionType},
		{"Save Wallpaper", fmt.Sprintf("%v", c.SaveWallpaper)},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Option", "Value"})
	table.AppendBulk(data)

	table.Render()
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

	tempDir := filepath.Join(dir, "temp")
	_, err = os.Stat(tempDir)
	if os.IsNotExist(err) {
		os.Mkdir(tempDir, os.ModePerm)
	}

	return dir
}

func jsonifyConfig(config *Configuration) []byte {
	byteArray, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		panic(err)
	}
	return (byteArray)
}

func ReadConfig() *Configuration {
	config := Configuration{}

	configFilePath := filepath.Join(getIrisDir(), "config.json")

	if !CheckFileExists(configFilePath) {
		defaultConfig := Configuration{
			SearchTerms:             []string{"nature"},
			Resolution:              "1600x900",
			ChangeWallpaper:         false,
			ChangeWallpaperDuration: -1,
			WallpaperDirectory:      "",
			SelectionType:           "random",
			SaveWallpaper:           false,
		}

		defaultConfig.WriteConfig()

		return &defaultConfig
	}

	configContent := readFile(configFilePath)
	if e := json.Unmarshal([]byte(configContent), &config); e != nil {
		panic(e)
	}

	return &config
}
