package internal

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// Update updates iris by downloading the latest executable from github, and renaming the
// old executable to `iris-old` so that it can be deleted by `DeletePreviousInstallation`.
func Update() {
	fmt.Println("Updating iris...")
	fmt.Println("Downloading the iris executable...")

	// * determining the os-specific url
	url := ""
	switch runtime.GOOS {
	case "windows":
		url = "https://github.com/Shravan-1908/iris/releases/latest/download/iris-windows-amd64.exe"
	// case "linux":
	// 	url = "https://github.com/Shravan-1908/iris/releases/latest/download/iris-linux-amd64"
	// case "darwin":
	// 	url = "https://github.com/Shravan-1908/iris/releases/latest/download/iris-darwin-amd64"
	default:
		fmt.Println("Your OS isn't supported by iris.")
		return
	}

	// * sending a request
	res, err := http.Get(url)

	if err != nil {
		fmt.Println("Error: Unable to download the executable. Check your internet connection.")
		fmt.Println(err.Error())
		return
	}

	defer res.Body.Close()

	// * determining the executable path
	downloadPath, e := os.UserHomeDir()
	if e != nil {
		fmt.Println("Error: Unable to determine iris's location.")
		fmt.Println(e.Error())
		return
	}
	downloadPath += "/.iris/iris"
	if runtime.GOOS == "windows" {
		downloadPath += ".exe"
	}

	os.Rename(downloadPath, downloadPath+"-old")

	exe, er := os.Create(downloadPath)
	if er != nil {
		fmt.Println("Error: Unable to access the file system.")
		fmt.Println(er.Error())
		return
	}
	defer exe.Close()

	// * writing the received content to the iris executable
	_, errr := io.Copy(exe, res.Body)
	if errr != nil {
		fmt.Println("Error: Unable to write the executable.")
		fmt.Println(errr.Error())
		return
	}

	// * performing an additional `chmod` utility for linux and mac
	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		os.Chmod(downloadPath, 0755)
	}

	fmt.Println("iris was updated successfully.")
}

// DeletePreviousInstallation deletes previous installation if it exists.
func DeletePreviousInstallation() {
	irisDir, _ := os.UserHomeDir()
	irisDir += "/.iris"

	files, _ := ioutil.ReadDir(irisDir)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "-old") {
			// fmt.Println("found existsing installation")
			os.Remove(irisDir + "/" + f.Name())
		}
		// fmt.Println(f.Name())
	}
}
