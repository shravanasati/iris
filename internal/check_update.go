package internal

import (
	"encoding/json"
	"fmt"
	"golang.org/x/mod/semver"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var lastCheckedFilePath = filepath.Join(GetIrisDir(), "last_checked.update")
var timeFormat = time.RFC3339

const (
	lower = iota
	equal
	greater
)

// Compares two semver strings of format vx.y.z.
// Returns greater if v1 > v2, lower if v1 < v2, equal if v1 == v2.
func compareSemverStrings(v1, v2 string) int {
	ans := semver.Compare(v1, v2)
	switch ans {
	case -1:
		return lower
	case 0:
		return equal
	case 1:
		return greater
	default:
		panic(fmt.Sprintf("unknown number returned by semver.Compare(%s, %s)=%d", v1, v2, ans))
	}
}

func getLastCheckedTime() time.Time {
	if !CheckPathExists(lastCheckedFilePath) {
		LogInfof("update", "last checked file not found, creating new one")
		lastCheckedFile, err := os.Create(lastCheckedFilePath)
		if err != nil {
			LogErrorf("update", "failed to create last checked file: %v", err)
			panic("unable to create the last checked file")
		}
		defer lastCheckedFile.Close()
		t := time.Now().Add(-1 * time.Hour * 24)
		writeLastCheckedTime(t)
		return t
	} else {
		content := readFile(lastCheckedFilePath)
		t, err := time.Parse(timeFormat, content)
		if err != nil {
			LogErrorf("update", "failed to parse last checked time: %v", err)
			panic("unable to parse time in last checked update file")
		}
		return t
	}
}

func writeLastCheckedTime(t time.Time) {
	LogInfof("update", "writing last checked time: %s", t.Format(timeFormat))
	lastCheckedFile, err := os.Create(lastCheckedFilePath)
	if err != nil {
		LogErrorf("update", "failed to write last checked time: %v", err)
		panic("unable to create the last checked update file")
	}
	lastCheckedFile.Write([]byte(t.Format(timeFormat)))
}

type releaseInfo struct {
	HTMLURL    string `json:"html_url"`
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	IsDraft    bool   `json:"draft"`
	IsPrelease bool   `json:"prerelease"`
	Body       string `json:"body"`
}

func (ri releaseInfo) display() {
	LogInfof("update", "notifying user of new release: %s", ri.TagName)
	fmt.Printf("\nnote: new `%s` of iris is now available at `%s` \n\nfull release article: \n%s \n%s\n", ri.TagName, ri.HTMLURL, ri.Name, ri.Body)
}

func CheckForUpdates(currentVersion string) {
	config := ReadConfig()
	if !config.CheckForUpdates {
		LogInfof("update", "update check disabled in config")
		return
	}

	now := time.Now()
	lastChecked := getLastCheckedTime()
	if !(now.Sub(lastChecked).Hours() >= 24) {
		// updates were checked for within 24 hours
		LogInfof("update", "skipping update check, last checked %v hours ago", now.Sub(lastChecked).Hours())
		return
	}

	LogInfof("update", "checking for updates, current version: %s", currentVersion)
	url := "https://api.github.com/repos/shravanasati/iris/releases/latest"
	releaseInfo := releaseInfo{}
	resp, err := http.Get(url)
	if err != nil {
		LogErrorf("update", "failed to fetch latest release: %v", err)
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		LogErrorf("update", "failed to read release response: %v", err)
		return
	}
	json.Unmarshal(data, &releaseInfo)
	writeLastCheckedTime(now)

	LogInfof("update", "latest release version: %s", releaseInfo.TagName)
	if compareSemverStrings(releaseInfo.TagName, currentVersion) != greater {
		// no new version
		LogInfof("update", "iris is up to date")
		return
	}
	if releaseInfo.IsDraft || releaseInfo.IsPrelease {
		// dont want to tell users about draft or prereleases
		LogInfof("update", "new release %s found but it is a draft or prerelease, skipping", releaseInfo.TagName)
		return
	}

	releaseInfo.display()
}
