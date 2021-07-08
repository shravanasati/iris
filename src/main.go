package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"errors"
	"github.com/reujab/wallpaper"
)

func stringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

func downloadImage(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("non-200 status code")
	}

	cacheDir := filepath.Join(getIrisDir(), "wallpapers")

	file, err := os.Create(filepath.Join(cacheDir, time.Now().Format("02-01-2006 15-04-05" + ".jpg")))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func main() {
	c := readConfig()

	useUnsplash := false
	if c.WallpaperDirectory == "" || !checkFileExists(c.WallpaperDirectory) {
		useUnsplash = true
	}

	resolution := c.Resolution
	if !stringInSlice(resolution, supportedResolutions) {
		resolution = "1600x900"
	}

	if useUnsplash {
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

	} else {
		err := wallpaper.SetFromFile(`C:\Users\LENOVO\Downloads\Images\wp9310707-summer-scotland-wallpapers.jpg`)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(c)
}
