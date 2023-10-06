package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Shravan-1908/go-wallpaper"
)

var validImageExtensions = []string{"png", "jpg", "jpeg", "jfif"}

// SetWallpaper sets the wallpaper to the given file.
func SetWallpaper(filename string) error {
	if !CheckFileExists(filename) {
		return fmt.Errorf("the file `%s` doesn't exist", filename)
	}

	absPath, err := (filepath.Abs(filename))
	if err != nil {
		return err
	}
	return wallpaper.SetFromFile(absPath)
}

// todo get wallpaper

func (c *Configuration) RemoteWallpaper() {
	switch strings.ToLower(strings.TrimSpace(c.RemoteSource)) {
	case "unsplash":
		c.unsplashWallpaper()
	case "spotlight":
		c.windowsSpotlightWallpaper()
	// todo match github url regex here
	case "github":
		c.githubRepoWallpaper()
	default:
		// todo edit readme about new config options - remote source and check for updates
		// todo link to remote source docs here
		fmt.Printf("Invalid remote source `%s`, defaulting to unsplash. Know more about iris configuration at https://github.com/Shravan-1908/iris#customization \n", c.RemoteSource)
		c.unsplashWallpaper()
	}
}

// unsplashWallpaper changes the wallpaper using unsplash.
func (c *Configuration) unsplashWallpaper() {
	searchTerms := strings.Join(c.SearchTerms, ",")

	url := fmt.Sprintf("https://source.unsplash.com/%v/?%v", c.Resolution, searchTerms)

	if !c.SaveWallpaper {
		f, e := downloadImage(url, true)
		if e != nil {
			fmt.Println(e)
		} else {
			if se := SetWallpaper(f); se != nil {
				fmt.Println(se.Error())
				os.Exit(1)
			}
		}
	} else {
		f, e := downloadImage(url, false)
		if e != nil {
			fmt.Println(e)
		} else {
			if se := SetWallpaper(f); se != nil {
				fmt.Println(se.Error())
				os.Exit(1)
			}
		}
	}
}

func (c *Configuration) windowsSpotlightWallpaper() {
	// searchTerms := strings.Join(c.SearchTerms, ",")
}

func (c *Configuration) githubRepoWallpaper() {}

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
		if StringInSlice(ext, validImageExtensions) {
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
	tempContents, er := ioutil.ReadDir(filepath.Join(GetIrisDir(), "temp"))
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
