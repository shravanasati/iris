/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// cacheCmd represents the cache command
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage iris cache data.",
	Long: `The current way video wallpapers are implemented in iris using ffmpeg, which
converts a given video into several thousand frames, and stores them in the cache folder.

The cache command thus provides the access to manage this cache.

Examples:
// todo add cache cmd examples here
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cache called")
	},
}

// todo document cache command in the readme
// todo ruminate over whether to remove a single item from the cache or not

func init() {
	rootCmd.AddCommand(cacheCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cacheCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cacheCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
