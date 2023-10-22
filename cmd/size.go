/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Shravan-1908/iris/internal"
	"github.com/spf13/cobra"
)

// sizeCmd represents the size command
var sizeCmd = &cobra.Command{
	Use:   "size",
	Short: "Prints the total cache size.",
	Long: `Prints the total cache size used by iris.

Example:
$ iris cache size`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Total cache size: ", internal.CacheSize())
	},
}

func init() {
	cacheCmd.AddCommand(sizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
