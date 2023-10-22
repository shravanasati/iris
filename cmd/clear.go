/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Shravan-1908/iris/internal"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the entire cache.",
	Long: `The clear command effectively deletes all caches and their references.

Example:
$ iris cache clear`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.CacheEmpty()
		if err != nil {
			fmt.Println("Unable to empty cache:", err)
		} else {
			fmt.Println("Successfully emptied cache.")
		}
	},
}

func init() {
	cacheCmd.AddCommand(clearCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
