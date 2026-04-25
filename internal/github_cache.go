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
		LogInfof("github-cache", "no cache file found at %s", cachePath)
		return cache
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		LogErrorf("github-cache", "failed to read cache file: %v", err)
		return cache
	}

	err = json.Unmarshal(data, &cache)
	if err != nil {
		LogErrorf("github-cache", "failed to unmarshal cache data: %v", err)
	}
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
		LogErrorf("github-cache", "failed to get github api url for %s: %v", repoURL, err)
		return err
	}

	// retrieve gh api token from config
	ghToken := c.GitHubAPIToken
	if ghToken == "" {
		ghToken = os.Getenv("IRIS_GH_TOKEN")
	}

	parsed, err := parseGitHubRepoSource(repoURL)
	if err != nil {
		LogErrorf("github-cache", "failed to parse github repo source %s: %v", repoURL, err)
		return err
	}

	LogInfof("github-cache", "syncing %s", repoURL)
	fmt.Printf("Syncing %s...\n", repoURL)
	_, err = FetchAndCache(parsed.normalizedURL, preparedURL, parsed.owner, parsed.repo, parsed.branch, parsed.folderPath, ghToken)
	if err != nil {
		LogErrorf("github-cache", "failed to sync %s: %v", repoURL, err)
	}
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
	LogInfof("github-cache", "starting fetch and cache for repo: %s", repoURL)

	cache := LoadGitHubCache()
	cachedData, exists := cache[repoURL]
	LogInfof("github-cache", "cache loaded. entries: %d, hit: %t", len(cache), exists)

	// 1. Try to get latest SHA
	shaStart := time.Now()
	LogInfof("github-cache", "checking latest commit sha")
	latestSHA, err := GetLatestCommitSHA(owner, repo, branch, folderPath, token)
	if err != nil {
		LogWarnf("github-cache", "failed to fetch latest commit sha: %v, falling back to cache if available", err)
		if exists {
			LogInfof("github-cache", "using cached wallpapers due to sha error. count: %d, elapsed: %s", len(cachedData.Wallpapers), time.Since(start))
			return cachedData.Wallpapers, nil
		}
		return nil, err
	}
	LogInfof("github-cache", "latest sha resolved in %s", time.Since(shaStart))

	// 2. Check if SHA matches
	if exists && cachedData.LastSHA == latestSHA {
		LogInfof("github-cache", "cache is up to date. returning cached wallpapers. count: %d, elapsed: %s", len(cachedData.Wallpapers), time.Since(start))
		return cachedData.Wallpapers, nil
	}

	// 3. Fetch fresh list
	fetchStart := time.Now()
	LogInfof("github-cache", "cache miss or stale sha. fetching fresh file list")
	req, err := http.NewRequest(http.MethodGet, preparedURL, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		LogWarnf("github-cache", "failed to fetch fresh file list: %v, falling back to cache if available", err)
		if exists {
			LogInfof("github-cache", "using cached wallpapers due to fetch error. count: %d, elapsed: %s", len(cachedData.Wallpapers), time.Since(start))
			return cachedData.Wallpapers, nil
		}
		return nil, err
	}
	defer resp.Body.Close()
	LogInfof("github-cache", "fresh list response status: %s in %s", resp.Status, time.Since(fetchStart))

	if resp.StatusCode != http.StatusOK {
		LogWarnf("github-cache", "github api returned %s, falling back to cache if available", resp.Status)
		if exists {
			LogInfof("github-cache", "using cached wallpapers due to non-200 response. count: %d, elapsed: %s", len(cachedData.Wallpapers), time.Since(start))
			return cachedData.Wallpapers, nil
		}
		return nil, fmt.Errorf("received non-200 status code: %s", resp.Status)
	}

	decodeStart := time.Now()
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		LogErrorf("github-cache", "failed to read response body: %v", err)
		return nil, err
	}

	var freshWallpapers []map[string]any
	if err := json.Unmarshal(jsonData, &freshWallpapers); err != nil {
		LogErrorf("github-cache", "failed to unmarshal fresh wallpapers: %v", err)
		return nil, err
	}
	LogInfof("github-cache", "decoded fresh wallpapers. count: %d in %s", len(freshWallpapers), time.Since(decodeStart))

	// 4. Update and Save Cache
	cache[repoURL] = GitHubRepoCache{
		LastSHA:    latestSHA,
		Wallpapers: freshWallpapers,
	}
	if err := SaveGitHubCache(cache); err != nil {
		LogErrorf("github-cache", "failed to save cache: %v", err)
		return nil, err
	}
	LogInfof("github-cache", "cache saved successfully. total elapsed: %s", time.Since(start))

	return freshWallpapers, nil
}
