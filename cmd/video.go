/*
Copyright © 2023 Shravan <dev.shravan@proton.me>
*/
package cmd

import (
	"fmt"

	"github.com/shravanasati/iris/internal"
	"github.com/spf13/cobra"
)

// videoCmd represents the video command
var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Set a video as wallpaper (experimental).",
	Long: `Set a video as wallpaper. This feature is currently experimental, and requires ffmpeg to work.
	
Example: iris video /path/to/video
You may want to start iris as a background process when setting a video wallpaper, suffix an ampersand to do the same. 
Example: iris video /path/to/video &
`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		if err := internal.SetVideoWallpaper(args[0]); err != nil {
			fmt.Println("Unable to set video wallpaper!")
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(videoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// videoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// videoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
