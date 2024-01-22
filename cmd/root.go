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
	"path/filepath"
	"strings"
	"time"

	"github.com/shravanasati/iris/internal"
	"github.com/spf13/cobra"
)

var c = internal.ReadConfig()

func realMain() {
	internal.ClearTemp()

	// * determining if to use remote source or local images
	useRemoteSource := false
	if strings.TrimSpace(c.WallpaperDirectory) == "" || !internal.CheckPathExists(c.WallpaperDirectory) {
		useRemoteSource = true
	}

	if strings.TrimSpace(c.SaveWallpaperDirectory) == "" || !internal.CheckPathExists(c.SaveWallpaperDirectory) {
		c.SaveWallpaperDirectory = filepath.Join(internal.GetIrisDir(), "wallpapers")
	}

	if !internal.ItemInSlice(c.Resolution, internal.SupportedResolutions) {
		c.Resolution = "1600x900"
	}

	// * wallpapers via remote source
	if useRemoteSource {
		if c.ChangeWallpaper {
			duration, e := time.ParseDuration(c.ChangeWallpaperDuration)

			if e != nil {
				duration = time.Minute * 5
			}
			for {
				c.RemoteWallpaper()
				time.Sleep(duration)
			}
		} else {
			c.RemoteWallpaper()
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
	
Visit https://github.com/shravanasati/iris for more information.`,

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
