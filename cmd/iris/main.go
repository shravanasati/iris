package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Shravan-1908/iris/internal"
	"github.com/thatisuday/commando"
)

const (
	NAME    string = "iris"
	VERSION string = "v0.2.0"
)

func main() {
	fmt.Println(NAME, VERSION)

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("The app has got a fatal error, and it cannot proceed further. \nPlease file a bug report at https://github.com/Shravan-1908/issues/new/choose, with the following error message. \n```\n%s\n```", err)
			os.Exit(1)
		}
	}()

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
						time.Sleep(time.Duration(duration) * (time.Minute))
						internal.ClearTemp()
					}
				} else {
					c.UnsplashWallpaper()
					internal.ClearTemp()
				}

				// * wallpapers via local directory
			} else {
				c.DirectoryWallpaper()
			}

		})

	commando.Parse(nil)

}
