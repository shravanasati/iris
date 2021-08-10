package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/reujab/wallpaper"
)

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
			if se := wallpaper.SetFromFile(f); se != nil {
				panic(se)
			}
		}
	} else {
		f, e := downloadImage(url, false)
		if e != nil {
			fmt.Println(e)
		} else {
			if se := wallpaper.SetFromFile(f); se != nil {
				panic(se)
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
		if strings.HasSuffix(f.Name(), ".png") || strings.HasSuffix(f.Name(), "jpg") || strings.HasSuffix(f.Name(), "jpeg") || strings.HasSuffix(f.Name(), "jfif") {
			contents = append(contents, filepath.Join(c.WallpaperDirectory, f.Name()))
		}
	}

	return contents
}

func (c *Configuration) DirectoryWallpaper() {
	contents := c.getValidWallpapers()

	if c.SelectionType == "random" {
		if c.ChangeWallpaper {
			if c.ChangeWallpaperDuration <= 0 {
				c.ChangeWallpaperDuration = 15
			}
			for {
				if err := wallpaper.SetFromFile(randomChoice(contents)); err != nil {
					panic(err)
				}
				time.Sleep(time.Duration(c.ChangeWallpaperDuration) * time.Minute)
			}
		} else {
			if err := wallpaper.SetFromFile(randomChoice(contents)); err != nil {
				panic(err)
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

				if err := wallpaper.SetFromFile(contents[i]); err != nil {
					panic(err)
				}

				time.Sleep(time.Duration(c.ChangeWallpaperDuration) * time.Minute)
			}

		} else {
			if err := wallpaper.SetFromFile(contents[0]); err != nil {
				panic(err)
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
