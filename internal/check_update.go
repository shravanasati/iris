package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var lastCheckedFilePath = filepath.Join(GetIrisDir(), "last_checked.update")
var timeFormat = time.RFC3339

func getLastCheckedTime() time.Time {
	if !CheckFileExists(lastCheckedFilePath) {
		lastCheckedFile, err := os.Create(lastCheckedFilePath)
		if err != nil {
			panic("unable to create the last checked file")
		}
		defer lastCheckedFile.Close()
		t := time.Now().Add(-1 * time.Hour * 24)
		writeLastCheckedTime(t)
		return t
	} else {
		t, err := time.Parse(timeFormat, readFile(lastCheckedFilePath))
		if err != nil {
			panic("unable to parse time in last checked update file")
		}
		return t
	}
}

func writeLastCheckedTime(t time.Time) {
	lastCheckedFile, err := os.Create(lastCheckedFilePath)
	if err != nil {
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
	fmt.Printf("\nnote: new `%s` of iris is now available at `%s` \n\nfull release article: \n%s \n%s\n", ri.TagName, ri.HTMLURL, ri.Name, ri.Body)
}

func CheckForUpdates(currentVersion string) {
	now := time.Now()
	if !(now.Sub(getLastCheckedTime()).Hours() >= 24) {
		// updates were checked for within 24 hours
		return
	}

	url := "https://api.github.com/repos/Shravan-1908/iris/releases/latest"
	releaseInfo := releaseInfo{}
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	json.Unmarshal(data, &releaseInfo)
	writeLastCheckedTime(now)

	if currentVersion == releaseInfo.TagName {
		// no new version
		return
	}
	if releaseInfo.IsDraft || releaseInfo.IsPrelease {
		// dont want to tell users about draft or prereleases
		return
	}

	releaseInfo.display()
}
