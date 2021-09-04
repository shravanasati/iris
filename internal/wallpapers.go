package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"runtime"

	"github.com/reujab/wallpaper"
)

// setWallpaper sets the wallpaper to the given file according to the OS.
func setWallpaper(filename string) error {
	switch runtime.GOOS {
	case "windows":
		return wallpaper.SetFromFile(filename)
	case "linux":
		cmd := exec.Command("feh", "--bg-scale", filename)
		return cmd.Run()
	case "darwin":
		cmd := exec.Command("osascript", "-e", "tell application \"Finder\" to set desktop picture to POSIX file " + filename)
		return cmd.Run()
	default:
		return fmt.Errorf("unsupported platform")
	}
}

// UnsplashWallpaper changes the wallpaper using unsplash.
func (c *Configuration) UnsplashWallpaper() {
	searchTerms := ""
	for i, v := range c.SearchTerms {
		if i == len(c.SearchTerms)-1 {
			searchTerms += v
		} else {
			searchTerms += v + ","
		}
	}

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
		ext := splitted[len(splitted) - 1]
		if StringInSlice(ext, []string{".png", ".jpg", ".jpeg", ".jfif"}) {
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
			if c.ChangeWallpaperDuration <= 0 {
				c.ChangeWallpaperDuration = 15
			}
			for {
				if err := setWallpaper(randomChoice(contents)); err != nil {
					fmt.Println("Unable to set wallpaper. Make sure you've `feh` installed if you're on a Linux system.")
					os.Exit(1)
				}
				time.Sleep(time.Duration(c.ChangeWallpaperDuration) * time.Minute)
			}
		} else {
			if err := setWallpaper(randomChoice(contents)); err != nil {
				fmt.Println("Unable to set wallpaper. Make sure you've `feh` installed if you're on a Linux system.")
				os.Exit(1)
			}
		}

	} else {
		if c.ChangeWallpaper {
			if c.ChangeWallpaperDuration <= 0 {
				c.ChangeWallpaperDuration = 15
			}
			wallpapers := c.getValidWallpapers()
			sort.Strings(wallpapers)
			for i := range wallpapers {
				if i == len(contents)-1 {
					i = 0
				}

				if err := setWallpaper(contents[i]); err != nil {
					fmt.Println("Unable to set wallpaper. Make sure you've `feh` installed if you're on a Linux system.")
					os.Exit(1)
				}

				time.Sleep(time.Duration(c.ChangeWallpaperDuration) * time.Minute)
			}

		} else {
			if err := setWallpaper(contents[0]); err != nil {
				fmt.Println("Unable to set wallpaper. Make sure you've `feh` installed if you're on a Linux system.")
				os.Exit(1)
			}
		}
	}
}

// ClearClutter deletes all the wallpapers present in ~/.iris/temp.
func ClearClutter() {
	tempContents, er := ioutil.ReadDir(filepath.Join(getIrisDir(), "temp"))
	if er != nil {
		fmt.Println(er)
		panic("unable to get ~/.iris/temp contents")
	}

	for _, f := range tempContents {
		if err := os.Remove(filepath.Join(getIrisDir(), "temp", f.Name())); err != nil {
			fmt.Println(err)
			panic("unable to delete " + f.Name())
		}
	}
}
