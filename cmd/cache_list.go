/*
Copyright © 2023 Shravan Asati <dev.shravan@proton.me>
MIT Licensed
*/
package cmd

import (
	"github.com/shravanasati/iris/internal"
	"github.com/spf13/cobra"
)

// cacheListCmd represents the cache list command
var cacheListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cached videos.",
	Long: `Lists out the paths of all videos iris has cached.

Example:
$ iris cache list`,
	Run: func(cmd *cobra.Command, args []string) {
		internal.CacheShow()
	},
}

func init() {
	cacheVideoCmd.AddCommand(cacheListCmd)
}
