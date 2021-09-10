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
	"github.com/Shravan-1908/iris/internal"
	"github.com/spf13/cobra"
)

var config = internal.ReadConfig()

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure iris for a personalized experience.",
	Long: `The config command is used to customize iris according to your needs. All configuration options are exposed as flags.
	
Example:
$ iris config -c 

`,
	Run: func(cmd *cobra.Command, args []string) {
		config.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().BoolVarP(&config.ChangeWallpaper, "change-wallpaper", "c", config.ChangeWallpaper, "Whether to change wallpapers continuosly in the background.")

	configCmd.Flags().IntVarP(&config.ChangeWallpaperDuration, "wallpaper-change-duration", "d", config.ChangeWallpaperDuration, "The duration between wallpaper changes, if to change them continuosly.")

	configCmd.Flags().BoolVarP(&config.SaveWallpaper, "save-wallpaper", "s", config.SaveWallpaper, "Whether to save the wallpaper to the local directory.")

	configCmd.Flags().StringVarP(&config.WallpaperDirectory, "wallpaper-directory", "w", config.WallpaperDirectory, "The local directory to get wallpapers from.")

	configCmd.Flags().StringVarP(&config.Resolution, "resolution", "r", config.Resolution, "The image resolution to use for unsplash wallpapers.")

	configCmd.Flags().StringVarP(&config.SelectionType, "selection-type", "t", config.SelectionType, "The selection type for choosing wallpapers from the local directory, either `random` or `sorted`.")

	configCmd.Flags().StringSliceVarP(&config.SearchTerms, "search-terms", "q", config.SearchTerms, "The search terms for unsplash wallpapers.")

}
