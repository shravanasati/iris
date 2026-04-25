package internal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/shravanasati/go-wallpaper"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

const (
	spotlightDomain = "https://windows10spotlight.com"
	searchEndpoint  = "/tag"
)

var (
	resolutionRegex = regexp.MustCompile(`-\d+x\d+`)
	protocolRegex   = regexp.MustCompile(`(?i)(http(s)*:(\/){2})`)

	// matches subreddits like r/sub1, r/sub1+sub2+sub3
	redditRegex = regexp.MustCompile(`^r/[\w\d_]{3,20}(?:\+[\w\d_]{3,20})*$`)

	// matches a remote github folder
	githubRegex          = regexp.MustCompile(`(?i)^((https:\/\/)*(github\.com))(\/[\w\-_\d% \.\(\)@&]+){2}\/tree(\/.+){1,}(\/){0,1}$`)
	getParamsGithubRegex = regexp.MustCompile(`(?i)^github\.com/([^/]+)/([^/]+)/tree/([^/]+)(/.*)?$`)
)

type githubRepoSource struct {
	normalizedURL string
	owner         string
	repo          string
	branch        string
	folderPath    string
}

var validImageExtensions = []string{"png", "jpg", "jpeg", "jfif"}

// SetWallpaper sets the wallpaper to thegiven file.
func SetWallpaper(filename string) error {
	if !CheckPathExists(filename) {
		return fmt.Errorf("the file `%s` doesn't exist", filename)
	}

	absPath, err := (filepath.Abs(filename))
	if err != nil {
		return err
	}
	return wallpaper.SetFromFile(absPath)
}

// Returns the current set wallpaper or the error.
func GetWallpaper() string {
	wallpaperPath, err := wallpaper.Get()
	if err != nil {
		return fmt.Sprintf("unable to get wallpaper: %v\n", err)
	}
	return wallpaperPath
}


// RemoteWallpaper dispatches the appropriate function to change wallpaper.
func (c *Configuration) RemoteWallpaper() {
	unquotedSource := strings.Trim(c.RemoteSource, "\"'")
	remoteSource := strings.ToLower(strings.TrimSpace(unquotedSource))
	if remoteSource == "spotlight" {
		if err := c.windowsSpotlightWallpaper(); err != nil {
			fmt.Println(err)
		}
	} else if redditRegex.Match([]byte(remoteSource)) {
		if err := c.redditWallpaper(); err != nil {
			fmt.Println(err)
		}
	} else if githubRegex.Match([]byte(remoteSource)) {
		if err := c.githubRepoWallpaper(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("Invalid remote source `%s`, defaulting to spotlight. Know more about iris remote source configuration at https://github.com/shravanasati/iris#customization \n", unquotedSource)
		if err := c.windowsSpotlightWallpaper(); err != nil {
			fmt.Println(err)
		}
	}
}


func (c *Configuration) windowsSpotlightWallpaper() error {
	// determine the url to hit
	var url string
	if len(c.SearchTerms) == 0 {
		url = spotlightDomain
	} else {
		searchTerms := strings.Join(c.SearchTerms, "+")
		url = spotlightDomain + searchEndpoint + "/" + searchTerms
	}

	// send a get request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("unable to load page: %s, error: %v", url, err)
	}
	defer resp.Body.Close()

	// parse the html content
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to parse html document from windows10spotlight: %v", err)
	}

	// find image links
	var links []string
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists && strings.Contains(src, "windows10spotlight") {
			link := resolutionRegex.ReplaceAllString(src, "")
			links = append(links, link)
		}
	})

	if len(links) == 0 {
		return fmt.Errorf("unable to find any image link on url=%v", url)
	}

	// select a random image, download it, and set it as wallpaper
	selectedURL := randomChoice(links)
	f, err := downloadImage(selectedURL, !c.SaveWallpaper)
	if err != nil {
		return fmt.Errorf("unable to download image: %s", selectedURL)
	}
	if err := SetWallpaper(f); err != nil {
		return fmt.Errorf("unable to set wallpaper: %s", err)
	}
	return nil
}

func getGithubAPIURL(ghRepoFolderURL string) (string, error) {
	parsed, err := parseGitHubRepoSource(ghRepoFolderURL)
	if err != nil {
		return "", err
	}
	preparedURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents%s?ref=%s", parsed.owner, parsed.repo, parsed.folderPath, parsed.branch)
	return preparedURL, nil
}

func parseGitHubRepoSource(ghRepoFolderURL string) (githubRepoSource, error) {
	if protocolRegex.Match([]byte(strings.ToLower(ghRepoFolderURL))) {
		ghRepoFolderURL = protocolRegex.ReplaceAllLiteralString(ghRepoFolderURL, "")
	}

	ghRepoFolderURL, _ = url.PathUnescape(ghRepoFolderURL)

	matches := getParamsGithubRegex.FindStringSubmatch(ghRepoFolderURL)
	var owner, repo, branch, folderPath string
	if len(matches) == 5 {
		owner = matches[1]
		repo = matches[2]
		branch = matches[3]
		folderPath = matches[4]
	} else if len(matches) == 4 {
		owner = matches[1]
		repo = matches[2]
		branch = matches[3]
		folderPath = ""
	} else {
		return githubRepoSource{}, fmt.Errorf("invalid remote source: %s. check your github URL, it must be of format github.com/owner/repo/tree/branch/optionalFolderPath", ghRepoFolderURL)
	}

	return githubRepoSource{
		normalizedURL: ghRepoFolderURL,
		owner:         owner,
		repo:          repo,
		branch:        branch,
		folderPath:    folderPath,
	}, nil
}

func (c *Configuration) githubRepoWallpaper() error {
	repoFolderURL := c.RemoteSource
	preparedURL, err := getGithubAPIURL(repoFolderURL)
	if err != nil {
		return err
	}

	// retrieve gh api token from config
	ghToken := c.GitHubAPIToken
	// if failed lookup environment variable
	if ghToken == "" {
		ghToken = os.Getenv("IRIS_GH_TOKEN")
	}

	parsed, err := parseGitHubRepoSource(repoFolderURL)
	if err != nil {
		return err
	}

	recvData, err := FetchAndCache(parsed.normalizedURL, preparedURL, parsed.owner, parsed.repo, parsed.branch, parsed.folderPath, ghToken)
	if err != nil {
		return err
	}

	// download image
	choice := randomChoice(recvData)["download_url"]
	downloadURL, ok := choice.(string)
	if !ok {
		return fmt.Errorf("unable to assert string type onto download url: %v", choice)
	}
	f, err := downloadImage(downloadURL, !c.SaveWallpaper)
	if err != nil {
		return err
	}

	// set downloaded image as wallpaper
	if err = SetWallpaper(f); err != nil {
		return err
	}
	return nil
}

func (c *Configuration) redditWallpaper() error {
	// todo add support for imgur, ireddit, gallery
	userAgent := fmt.Sprintf("%v:iris-%v:v0.4.0 (by /u/%v)", runtime.GOOS, _UUID[:6], _UUID[:6])
	client, err := reddit.NewReadonlyClient(reddit.WithUserAgent(userAgent))
	if err != nil {
		return err
	}
	// todo use reddit token if found
	subredditName := strings.Replace(strings.ToLower(c.RemoteSource), "r/", "", 1)
	posts, _, err := client.Subreddit.TopPosts(context.Background(), subredditName, &reddit.ListPostOptions{Time: "all"})
	if err != nil {
		return err
	}
	f, err := downloadImage(randomChoice(posts).URL, !c.SaveWallpaper)
	if err != nil {
		return err
	}
	err = SetWallpaper(f)
	if err != nil {
		return err
	}
	// todo how to download gallery posts
	// todo match reddit similar to github, r/wallpapers/top?t=all&limit=50

	return nil
}


func (c *Configuration) getValidWallpapers() []string {
	contents := []string{}
	tempContents, er := os.ReadDir(c.WallpaperDirectory)
	if er != nil {
		panic(er)
	}

	for _, f := range tempContents {
		splitted := strings.Split(f.Name(), ".")
		if len(splitted) == 0 {
			continue
		}
		ext := strings.ToLower(splitted[len(splitted)-1])
		if ItemInSlice(ext, validImageExtensions) {
			contents = append(contents, filepath.Join(c.WallpaperDirectory, f.Name()))
		}
	}

	return contents
}

func (c *Configuration) DirectoryWallpaper() {
	contents := c.getValidWallpapers()
	if len(contents) == 0 {
		fmt.Printf("No valid wallpapers found in the directory `%s`.\n", c.WallpaperDirectory)
		return
	}

	if c.SelectionType == "random" {
		if c.ChangeWallpaper {
			duration, e := time.ParseDuration(c.ChangeWallpaperDuration)
			if e != nil {
				duration = time.Minute * 5
			}
			for {
				if err := SetWallpaper(randomChoice(contents)); err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				time.Sleep(duration)
			}
		} else {
			if err := SetWallpaper(randomChoice(contents)); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

	} else {
		if c.ChangeWallpaper {
			duration, e := time.ParseDuration(c.ChangeWallpaperDuration)
			if e != nil {
				duration = time.Minute * 5
			}

			wallpapers := c.getValidWallpapers()
			sort.Strings(wallpapers)
			for {
				for i := range wallpapers {
					if err := SetWallpaper(contents[i]); err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}

					time.Sleep(duration)
				}
			}

		} else {
			if err := SetWallpaper(contents[0]); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}
	}
}

// ClearTemp deletes all the wallpapers present in ~/.iris/temp.
func ClearTemp() {
	tempContents, er := os.ReadDir(filepath.Join(GetIrisDir(), "temp"))
	if er != nil {
		fmt.Println(er)
		panic("unable to get ~/.iris/temp contents")
	}

	for _, f := range tempContents {
		if err := os.Remove(filepath.Join(GetIrisDir(), "temp", f.Name())); err != nil {
			fmt.Println(err)
			panic("unable to delete " + f.Name())
		}
	}
}
