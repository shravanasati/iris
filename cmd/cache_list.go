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
	cacheCmd.AddCommand(cacheListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cacheListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cacheListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
