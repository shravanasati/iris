package internal

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"os"

	"github.com/reujab/wallpaper"
)

func UnsplashWallpaper(c *Configuration, resolution string) {
	searchTerms := ""
	for i, v := range c.SearchTerms {
		if i == len(c.SearchTerms)-1 {
			searchTerms += v
		} else {
			searchTerms += v + ","
		}
	}

	url := fmt.Sprintf("https://source.unsplash.com/%v/?%v", resolution, searchTerms)

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

func getValidWallpapers(c *Configuration) []string {
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

func DirectoryWallpaper(c *Configuration) {
	contents := getValidWallpapers(c)

	if c.SelectionType == "random" {
		err := wallpaper.SetFromFile(randomChoice(contents))
		if err != nil {
			panic(err)
		}
	} else {
		err := wallpaper.SetFromFile(contents[0])
		if err != nil {
			panic(err)
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