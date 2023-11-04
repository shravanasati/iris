package internal

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Shravan-1908/go-wallpaper"
)

const (
	spotlightDomain = "https://windows10spotlight.com"
	searchEndpoint  = "/tag"
)

var (
	// matches a remote github folder
	githubRegex = regexp.MustCompile(`(?i)^((https:\/\/)*(github\.com))(\/[\w\-_\d]+){2}\/tree(\/[\w\-_\d]+){1,}$`)
)

var validImageExtensions = []string{"png", "jpg", "jpeg", "jfif"}

// SetWallpaper sets the wallpaper to thegiven file.
func SetWallpaper(filename string) error {
	if !CheckPathExists(filename) {
		return fmt.Errorf("the file `%s` doesn't exist", filename)
	}

	absPath, err := (filepath.Abs(filename))
	if err != nil {
		return err
	}
	return wallpaper.SetFromFile(absPath)
}

// Returns the current set wallpaper or the error.
func GetWallpaper() string {
	wallpaperPath, err := wallpaper.Get()
	if err != nil {
		return fmt.Sprintf("unable to get wallpaper: %v\n", err)
	}
	return wallpaperPath
}

// RemoteWallpaper dispatches the appropriate function to change wallpaper.
func (c *Configuration) RemoteWallpaper() {
	remoteSource := strings.ToLower(strings.TrimSpace(c.RemoteSource)) 
	if remoteSource == "unsplash" {
		c.unsplashWallpaper()
	} else if remoteSource == "spotlight" {
		c.windowsSpotlightWallpaper()
	} else if githubRegex.Match([]byte(remoteSource)) {
		c.githubRepoWallpaper()
	} else {
		// todo edit readme about new config options - remote source and check for updates
		// todo link to remote source docs here
		fmt.Printf("Invalid remote source `%s`, defaulting to unsplash. Know more about iris remote source configuration at https://github.com/Shravan-1908/iris#customization \n", c.RemoteSource)
		c.unsplashWallpaper()
	}
}

// unsplashWallpaper changes the wallpaper using unsplash.
func (c *Configuration) unsplashWallpaper() {
	searchTerms := strings.Join(c.SearchTerms, ",")

	url := fmt.Sprintf("https://source.unsplash.com/%v/?%v", c.Resolution, searchTerms)
	f, e := downloadImage(url, !c.SaveWallpaper)
	if e != nil {
		fmt.Println(e)
	} else {
		if se := SetWallpaper(f); se != nil {
			fmt.Println(se.Error())
			os.Exit(1)
		}
	}
}

func (c *Configuration) windowsSpotlightWallpaper() {
	searchTerms := strings.Join(c.SearchTerms, "+")
	url := spotlightDomain + searchEndpoint + "/" + searchTerms
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Unable to load page:", url)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Unable to parse html document from windows10spotlight")
		return
	}
	var links []string
	resolutionRegex := regexp.MustCompile(`-\d+x\d+`)
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && strings.Contains(src, "windows10spotlight") {
			link := resolutionRegex.ReplaceAllString(src, "")
			links = append(links, link)
		}
	})

	selectedURL := randomChoice(links)
	f, err := downloadImage(selectedURL, !c.SaveWallpaper)
	if err != nil {
		fmt.Println("Unable to download image:", selectedURL)
		return
	}
	if err := SetWallpaper(f); err != nil {
		fmt.Println("Unable to set wallpaper:", err)
		os.Exit(1)
	}
}

func (c *Configuration) githubRepoWallpaper() {
	fmt.Println("github repo wallpaper")
}

func (c *Configuration) getValidWallpapers() []string {
	contents := []string{}
	tempContents, er := os.ReadDir(c.WallpaperDirectory)
	if er != nil {
		panic(er)
	}

	for _, f := range tempContents {
		splitted := strings.Split(f.Name(), ".")
		if len(splitted) == 0 {
			continue
		}
		ext := strings.ToLower(splitted[len(splitted)-1])
		if ItemInSlice(ext, validImageExtensions) {
			contents = append(contents, filepath.Join(c.WallpaperDirectory, f.Name()))
		}
	}

	return contents
}

func (c *Configuration) DirectoryWallpaper() {
	contents := c.getValidWallpapers()
	if len(contents) == 0 {
		fmt.Printf("No valid wallpapers found in the directory `%s`.\n", c.WallpaperDirectory)
		return
	}

	if c.SelectionType == "random" {
		if c.ChangeWallpaper {
			duration, e := time.ParseDuration(c.ChangeWallpaperDuration)
			if e != nil {
				duration = time.Minute * 5
			}
			for {
				if err := SetWallpaper(randomChoice(contents)); err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				time.Sleep(duration)
			}
		} else {
			if err := SetWallpaper(randomChoice(contents)); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

	} else {
		if c.ChangeWallpaper {
			duration, e := time.ParseDuration(c.ChangeWallpaperDuration)
			if e != nil {
				duration = time.Minute * 5
			}

			wallpapers := c.getValidWallpapers()
			sort.Strings(wallpapers)
			for {
				for i := range wallpapers {
					if err := SetWallpaper(contents[i]); err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}

					time.Sleep(duration)
				}
			}

		} else {
			if err := SetWallpaper(contents[0]); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	}
}

// ClearTemp deletes all the wallpapers present in ~/.iris/temp.
func ClearTemp() {
	tempContents, er := os.ReadDir(filepath.Join(GetIrisDir(), "temp"))
	if er != nil {
		fmt.Println(er)
		panic("unable to get ~/.iris/temp contents")
	}

	for _, f := range tempContents {
		if err := os.Remove(filepath.Join(GetIrisDir(), "temp", f.Name())); err != nil {
			fmt.Println(err)
			panic("unable to delete " + f.Name())
		}
	}
}
