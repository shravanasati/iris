package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shravan-1908/iris/internal"
	"github.com/thatisuday/commando"
)

const (
	NAME    string = "iris"
	VERSION string = "v0.1.1"
)

func main() {
	fmt.Println(NAME, VERSION)
	go internal.DeletePreviousInstallation()

	commando.
		SetExecutableName(NAME).
		SetVersion(VERSION).
		SetDescription("iris is an easy to use and customizable wallpaper manager for windows.")

	// root command
	commando.
		Register(nil).
		SetShortDescription("Run iris").
		SetDescription("Runs iris.").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {

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
				c.DirectoryWallpaper()
			}

		})

	// update command
	commando.
		Register("up").
		SetShortDescription("Update iris").
		SetDescription("Updates iris to the latest version.").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			internal.Update()
		})

	commando.Parse(nil)

}
