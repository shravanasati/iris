package cmd

import (
	"github.com/spf13/cobra"
)

var cacheVideoCmd = &cobra.Command{
	Use:   "video",
	Short: "Manage video wallpaper cache",
}

func init() {
	cacheCmd.AddCommand(cacheVideoCmd)
}
