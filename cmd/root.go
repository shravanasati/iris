/*
Copyright © 2021 Shravan Asati <dev.shravan@protonmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"strings"
	"time"

	"github.com/Shravan-1908/iris/internal"
	"github.com/spf13/cobra"
)

var c = internal.ReadConfig()

func realMain() {
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
}

var rootCmd = &cobra.Command{
	Use:   "iris",
	Short: "Run iris.",
	Long: `iris is an easy to use, cross platform, feature rich, customizable and open source wallpaper manager. 
	
Visit https://github.com/Shravan-1908/iris for more information.`,

	Run: func(cmd *cobra.Command, args []string) {
		realMain()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.MousetrapHelpText = ""
}
