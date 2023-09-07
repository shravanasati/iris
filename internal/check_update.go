package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type releaseInfo struct {
	HTMLURL    string `json:"html_url"`
	TagName    string `json:"tag_name"`
	Name       string `json:"name"`
	IsDraft    bool   `json:"draft"`
	IsPrelease bool   `json:"prerelease"`
	Body       string
}

func (ri releaseInfo) display() {
	fmt.Printf("\nnote: new `%s` of iris is now available at `%s` \n\nfull release article: \ntitle: %s \nbody: %s\n", ri.TagName, ri.HTMLURL, ri.Name, ri.Body)
}

// todo check for updates only periodically

func CheckForUpdates(currentVersion string) {
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
