package internal

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

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
		err := wallpaper.SetFromURL(url)
		if err != nil {
			panic(err)
		}
	} else {
		f, e := downloadImage(url)
		if e != nil {
			fmt.Println(e)
		} else {
			wallpaper.SetFromFile(f)
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
