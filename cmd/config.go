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
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure iris for a personalized experience.",
	Long: `The config command is used to customize iris according to your needs. All configuration options are exposed as flags.
	
Examples:

$ iris config --remote-source spotlight
$ iris config --search-terms landscape,nature
$ iris config --save-wallpaper[=false]
$ iris config --wallpaper-directory /home/user/Pictures/Wallpapers
$ iris config --change-wallpaper[=false]
$ iris config --resolution 1920x1080
$ iris config list

`,
	Run: func(cmd *cobra.Command, args []string) {
		c.WriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVar(&c.RemoteSource, "remote-source", c.RemoteSource, "Remote source to select wallpapers from. Valid options are: unsplash, spotlight, github.")

	configCmd.Flags().BoolVar(&c.CheckForUpdates, "check-for-updates", c.CheckForUpdates, "Whether to check for updates of iris from github.")

	configCmd.Flags().BoolVarP(&c.ChangeWallpaper, "change-wallpaper", "c", c.ChangeWallpaper, "Whether to change wallpapers continuosly in the background.")

	configCmd.Flags().StringVarP(&c.ChangeWallpaperDuration, "wallpaper-change-duration", "d", c.ChangeWallpaperDuration, "The duration between wallpaper changes, if to change them continuosly.")

	configCmd.Flags().BoolVarP(&c.SaveWallpaper, "save-wallpaper", "s", c.SaveWallpaper, "Whether to save the wallpaper to the local directory.")

	configCmd.Flags().StringVarP(&c.WallpaperFile, "wallpaper-file", "f", c.WallpaperFile, "Path to the wallpaper file.")

	configCmd.Flags().StringVarP(&c.WallpaperDirectory, "wallpaper-directory", "w", c.WallpaperDirectory, "The local directory to get wallpapers from.")

	configCmd.Flags().StringVarP(&c.SaveWallpaperDirectory, "save-wallpaper-directory", "u", c.SaveWallpaperDirectory, "The local directory to save wallpapers in.")

	configCmd.Flags().StringVarP(&c.Resolution, "resolution", "r", c.Resolution, "The image resolution to use for unsplash wallpapers.")

	configCmd.Flags().StringVarP(&c.SelectionType, "selection-type", "t", c.SelectionType, "The selection type for choosing wallpapers from the local directory, either `random` or `sorted`.")

	configCmd.Flags().StringSliceVarP(&c.SearchTerms, "search-terms", "q", c.SearchTerms, "The search terms for unsplash wallpapers.")

	configCmd.Flags().StringVar(&c.GitHubAPIToken, "github-token", c.GitHubAPIToken, "The GitHub API token, used to perform authorized requests.")

}
