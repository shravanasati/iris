/*
Copyright © 2023 Shravan Asati <dev.shravan@proton.me>
MIT Licensed
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/shravanasati/iris/internal"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm <video-path>",
	Short: "Remove a specific video from the cache.",
	Long: `Removes a specific video and its associated frames from the iris cache.

Example:
$ iris cache rm "C:\Path\To\Video.mp4"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		videoPath := args[0]
		err := internal.CacheRemove(videoPath)
		if err != nil {
			fmt.Println("Error removing cache item:", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully removed %s from cache.\n", videoPath)
	},
}

func init() {
	cacheCmd.AddCommand(rmCmd)
}
