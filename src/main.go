package main

import (
	"fmt"
	"time"
)

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
		if c.ChangeWallpaper {
			duration := c.ChangeWallpaperDuration
			if duration <= 0 {
				duration = 5
			}
			for {
				unsplashWallpaper(c, resolution)
				time.Sleep(time.Duration(duration * int(time.Minute)))
			}
		} else {
			unsplashWallpaper(c, resolution)
		}

	} else {
		if c.ChangeWallpaper {
			duration := c.ChangeWallpaperDuration
			if duration <= 0 {
				duration = 5
			}
			for {
				directoryWallpaper(c)
				time.Sleep(time.Duration(duration * int(time.Minute)))
			}
		} else {
			directoryWallpaper(c)
		}
	}

	fmt.Println(c)
}
