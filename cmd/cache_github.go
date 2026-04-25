package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/shravanasati/iris/internal"
	"github.com/spf13/cobra"
)

var cacheGithubCmd = &cobra.Command{
	Use:   "github",
	Short: "Manage GitHub wallpaper cache",
}

var cacheGithubListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cached GitHub repositories.",
	Run: func(cmd *cobra.Command, args []string) {
		internal.GitHubCacheShow()
	},
}

var cacheGithubClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all cached GitHub repository results.",
	Run: func(cmd *cobra.Command, args []string) {
		cache := make(internal.GitHubCacheMap)
		if err := internal.SaveGitHubCache(cache); err != nil {
			fmt.Printf("Error clearing GitHub cache: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("GitHub cache cleared successfully.")
	},
}

var cacheGithubSyncCmd = &cobra.Command{
	Use:   "sync [repo-url]",
	Short: "Force sync the GitHub cache for a specific repository.",
	Long:  `Forces a refresh of the GitHub cache for the specified repository URL or the current one in config.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := internal.ReadConfig()
		repoURL := config.RemoteSource
		if len(args) > 0 {
			repoURL = args[0]
		}

		if !strings.Contains(repoURL, "github.com") {
			fmt.Println("Error: Invalid GitHub URL.")
			os.Exit(1)
		}

		if err := config.SyncGitHubCache(repoURL); err != nil {
			fmt.Printf("Error syncing cache: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully synced cache for %s.\n", repoURL)
	},
}

func init() {
	cacheCmd.AddCommand(cacheGithubCmd)
	cacheGithubCmd.AddCommand(cacheGithubListCmd)
	cacheGithubCmd.AddCommand(cacheGithubClearCmd)
	cacheGithubCmd.AddCommand(cacheGithubSyncCmd)
}
