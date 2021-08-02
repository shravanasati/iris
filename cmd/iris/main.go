package main

import (
	"fmt"
	"time"
	"github.com/Shravan-1908/iris/internal"
)

func main() {
	c := internal.ReadConfig()

	useUnsplash := false
	if c.WallpaperDirectory == "" || !internal.CheckFileExists(c.WallpaperDirectory) {
		useUnsplash = true
	}

	resolution := c.Resolution
	if !internal.StringInSlice(resolution, internal.SupportedResolutions) {
		resolution = "1600x900"
	}

	if useUnsplash {
		if c.ChangeWallpaper {
			duration := c.ChangeWallpaperDuration
			if duration <= 0 {
				duration = 5
			}
			for {
				internal.UnsplashWallpaper(c, resolution)
				time.Sleep(time.Duration(duration * int(time.Minute)))
			}
		} else {
			internal.UnsplashWallpaper(c, resolution)
		}

	} else {
		if c.ChangeWallpaper {
			duration := c.ChangeWallpaperDuration
			if duration <= 0 {
				duration = 5
			}
			for {
				internal.DirectoryWallpaper(c)
				time.Sleep(time.Duration(duration * int(time.Minute)))
			}
		} else {
			internal.DirectoryWallpaper(c)
		}
	}

	fmt.Println(c)
	internal.ClearClutter()
}
