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

// setWallpaper sets the wallpaper to the given file according to the OS.
func setWallpaper(filename string) error {
	absPath, err := (filepath.Abs(filename))
	if err != nil {
		return err
	}
	fmt.Println(absPath)
	return wallpaper.SetFromFile(absPath)
}

// UnsplashWallpaper changes the wallpaper using unsplash.
func (c *Configuration) UnsplashWallpaper() {
	searchTerms := strings.Join(c.SearchTerms, ",")

	url := fmt.Sprintf("https://source.unsplash.com/%v/?%v", c.Resolution, searchTerms)

	if !c.SaveWallpaper {
		f, e := downloadImage(url, true)
		if e != nil {
			fmt.Println(e)
		} else {
			if se := setWallpaper(f); se != nil {
				fmt.Println("Unable to set wallpaper. Make sure you've `feh` installed if you're on a Linux system.")
				os.Exit(1)
			}
		}
	} else {
		f, e := downloadImage(url, false)
		if e != nil {
			fmt.Println(e)
		} else {
			if se := setWallpaper(f); se != nil {
				fmt.Println("Unable to set wallpaper. Make sure you've `feh` installed if you're on a Linux system.")
				os.Exit(1)
			}
		}
	}
}

func (c *Configuration) getValidWallpapers() []string {
	contents := []string{}
	tempContents, er := ioutil.ReadDir(c.WallpaperDirectory)
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
				if err := setWallpaper(randomChoice(contents)); err != nil {
					fmt.Println("Unable to set wallpaper. Make sure you've `feh` installed if you're on a Linux system.")
					os.Exit(1)
				}
				time.Sleep(duration)
			}
		} else {
			if err := setWallpaper(randomChoice(contents)); err != nil {
				fmt.Println("Unable to set wallpaper. Make sure you've `feh` installed if you're on a Linux system.")
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
					if err := setWallpaper(contents[i]); err != nil {
						fmt.Println("Unable to set wallpaper. Make sure you've `feh` installed if you're on a Linux system.")
						os.Exit(1)
					}

					time.Sleep(duration)
				}
			}

		} else {
			if err := setWallpaper(contents[0]); err != nil {
				fmt.Println("Unable to set wallpaper. Make sure you've `feh` installed if you're on a Linux system.")
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
