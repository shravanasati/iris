package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type GitHubRepoCache struct {
	LastSHA    string           `json:"last_sha"`
	Wallpapers []map[string]any `json:"wallpapers"`
}

type GitHubCacheMap map[string]GitHubRepoCache

const githubAPITimeout = 20 * time.Second

func getGitHubCachePath() string {
	return filepath.Join(GetIrisDir(), "cache", "github.json")
}

func LoadGitHubCache() GitHubCacheMap {
	cachePath := getGitHubCachePath()
	cache := make(GitHubCacheMap)

	if !CheckPathExists(cachePath) {
		return cache
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		return cache
	}

	json.Unmarshal(data, &cache)
	return cache
}

func SaveGitHubCache(cache GitHubCacheMap) error {
	cachePath := getGitHubCachePath()
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}

func (c *Configuration) SyncGitHubCache(repoURL string) error {
	preparedURL, err := getGithubAPIURL(repoURL)
	if err != nil {
		return err
	}

	// retrieve gh api token from config
	ghToken := c.GitHubAPIToken
	if ghToken == "" {
		ghToken = os.Getenv("IRIS_GH_TOKEN")
	}

	parsed, err := parseGitHubRepoSource(repoURL)
	if err != nil {
		return err
	}

	fmt.Printf("Syncing %s...\n", repoURL)
	_, err = FetchAndCache(parsed.normalizedURL, preparedURL, parsed.owner, parsed.repo, parsed.branch, parsed.folderPath, ghToken)
	return err
}

func GetLatestCommitSHA(owner, repo, branch, folderPath, token string) (string, error) {
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)
	params := url.Values{}
	params.Set("sha", branch)
	params.Set("path", strings.TrimPrefix(folderPath, "/"))
	params.Set("per_page", "1")
	apiURL := baseURL + "?" + params.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), githubAPITimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return "", err
	}

	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	client := &http.Client{Timeout: githubAPITimeout}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timed out after %s while checking latest commit SHA", githubAPITimeout)
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %s", resp.Status)
	}

	var commits []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return "", err
	}

	if len(commits) == 0 {
		return "", fmt.Errorf("no commits found for path %s", folderPath)
	}

	sha, ok := commits[0]["sha"].(string)
	if !ok {
		return "", fmt.Errorf("unable to parse SHA from response")
	}

	return sha, nil
}

func GitHubCacheShow() {
	cache := LoadGitHubCache()
	if len(cache) == 0 {
		fmt.Println("No GitHub repositories cached.")
		return
	}

	tableData := [][]string{}
	i := 1
	for repo, data := range cache {
		tableData = append(tableData, []string{
			fmt.Sprintf("%v", i),
			repo,
			data.LastSHA[:7],
			fmt.Sprintf("%d", len(data.Wallpapers)),
		})
		i++
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"S. No.", "Repository", "SHA", "Wallpapers"})
	table.AppendBulk(tableData)
	table.Render()
}

func FetchAndCache(repoURL, preparedURL, owner, repo, branch, folderPath, token string) ([]map[string]any, error) {
	start := time.Now()
	fmt.Printf("[github-cache] start repo=%s\n", repoURL)

	cache := LoadGitHubCache()
	cachedData, exists := cache[repoURL]
	fmt.Printf("[github-cache] cache loaded entries=%d hit=%t\n", len(cache), exists)

	// 1. Try to get latest SHA
	shaStart := time.Now()
	fmt.Printf("[github-cache] checking latest commit SHA...\n")
	latestSHA, err := GetLatestCommitSHA(owner, repo, branch, folderPath, token)
	if err != nil {
		fmt.Printf("Warning: Failed to fetch latest commit SHA: %v. Falling back to cache if available.\n", err)
		if exists {
			fmt.Printf("[github-cache] using cached wallpapers due to SHA error count=%d elapsed=%s\n", len(cachedData.Wallpapers), time.Since(start))
			return cachedData.Wallpapers, nil
		}
		return nil, err
	}
	fmt.Printf("[github-cache] latest SHA resolved in %s\n", time.Since(shaStart))

	// 2. Check if SHA matches
	if exists && cachedData.LastSHA == latestSHA {
		fmt.Printf("[github-cache] cache is up to date; returning cached wallpapers count=%d elapsed=%s\n", len(cachedData.Wallpapers), time.Since(start))
		return cachedData.Wallpapers, nil
	}

	// 3. Fetch fresh list
	fetchStart := time.Now()
	fmt.Printf("[github-cache] cache miss or stale SHA; fetching fresh file list...\n")
	req, err := http.NewRequest(http.MethodGet, preparedURL, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Warning: Failed to fetch fresh file list: %v. Falling back to cache if available.\n", err)
		if exists {
			fmt.Printf("[github-cache] using cached wallpapers due to fetch error count=%d elapsed=%s\n", len(cachedData.Wallpapers), time.Since(start))
			return cachedData.Wallpapers, nil
		}
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Printf("[github-cache] fresh list response status=%s in %s\n", resp.Status, time.Since(fetchStart))

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Warning: GitHub API returned %s. Falling back to cache if available.\n", resp.Status)
		if exists {
			fmt.Printf("[github-cache] using cached wallpapers due to non-200 response count=%d elapsed=%s\n", len(cachedData.Wallpapers), time.Since(start))
			return cachedData.Wallpapers, nil
		}
		return nil, fmt.Errorf("received non-200 status code: %s", resp.Status)
	}

	decodeStart := time.Now()
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var freshWallpapers []map[string]any
	if err := json.Unmarshal(jsonData, &freshWallpapers); err != nil {
		return nil, err
	}
	fmt.Printf("[github-cache] decoded fresh wallpapers count=%d in %s\n", len(freshWallpapers), time.Since(decodeStart))

	// 4. Update and Save Cache
	cache[repoURL] = GitHubRepoCache{
		LastSHA:    latestSHA,
		Wallpapers: freshWallpapers,
	}
	if err := SaveGitHubCache(cache); err != nil {
		return nil, err
	}
	fmt.Printf("[github-cache] cache saved successfully total elapsed=%s\n", time.Since(start))

	return freshWallpapers, nil
}
