package internal

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func StringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

func randomChoice(slice []string) string {
	rand.Seed(time.Now().UnixNano())

	return slice[rand.Intn(len(slice))]
}

func CheckFileExists(filePath string) bool {
	_, e := os.Stat(filePath)
	return !os.IsNotExist(e)
}

func downloadImage(url string, temp bool) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("non-200 status code")
	}

	tempDir := filepath.Join(getIrisDir(), "wallpapers")
	if temp {
		tempDir = filepath.Join(getIrisDir(), "temp")
	}

	file, err := os.Create(filepath.Join(tempDir, time.Now().Format("02-01-2006 15-04-05"+".jpg")))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	err = file.Close()
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
