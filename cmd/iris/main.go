package main

import (
	"github.com/Shravan-1908/iris/internal"
	"strings"
	"time"
)

func main() {
	// * getting the configuration
	c := internal.ReadConfig()

	// * determining if to use unsplash or local images
	useUnsplash := false
	if strings.TrimSpace(c.WallpaperDirectory) == "" || !internal.CheckFileExists(c.WallpaperDirectory) {
		useUnsplash = true
	}

	resolution := c.Resolution
	if !internal.StringInSlice(resolution, internal.SupportedResolutions) {
		c.Resolution = "1600x900"
	}

	// * wallpapers via unsplash
	if useUnsplash {
		if c.ChangeWallpaper {
			duration := c.ChangeWallpaperDuration
			if duration <= 0 {
				duration = 15
			}
			for {
				c.UnsplashWallpaper()
				time.Sleep(time.Duration(duration * int(time.Minute)))
				internal.ClearClutter()
			}
		} else {
			c.UnsplashWallpaper()
			internal.ClearClutter()
		}

	// * wallpapers via local directory
	} else {
		// if c.ChangeWallpaper {
		// 	duration := c.ChangeWallpaperDuration
		// 	if duration <= 0 {
		// 		duration = 5
		// 	}
		// 	for {
		// 		internal.DirectoryWallpaper(c)
		// 		time.Sleep(time.Duration(duration * int(time.Minute)))
		// 	}
		// } else {
		// 	internal.DirectoryWallpaper(c)
		// }
		c.DirectoryWallpaper()
	}

}
