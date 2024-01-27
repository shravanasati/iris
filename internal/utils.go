package internal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

// Checks whether the given item exists in the slice.
func ItemInSlice[T comparable](s T, slice []T) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}

	return false
}

// Returns a random element from the given slice.
func randomChoice[T any](slice []T) T {
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))

	return slice[randGen.Intn(len(slice))]
}

// Checks whether a given path exists.
func CheckPathExists(filePath string) bool {
	_, e := os.Stat(filePath)
	return !os.IsNotExist(e)
}

// Downloads the image from the given URL. `temp` parameter is used to determine where to save the
// downloaded image.
// Returns filepath to the downloaded image and a error, if any.
func downloadImage(url string, temp bool) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("non-200 status code")
	}

	tempDir := ReadConfig().SaveWallpaperDirectory
	if !CheckPathExists(tempDir) {
		tempDir = filepath.Join(GetIrisDir(), "wallpapers")
	}
	if temp {
		tempDir = filepath.Join(GetIrisDir(), "temp")
	}

	filename := time.Now().Format("02-01-2006 15-04-05" + ".jpg")
	filename = strings.ReplaceAll(filename, " ", "-")
	file, err := os.Create(filepath.Join(tempDir, filename))

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

// readFile reads the given file and returns the string content of the same.
func readFile(file string) string {
	f, ferr := os.Open(file)
	if ferr != nil {
		panic(ferr)
	}
	defer f.Close()

	text := ""
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text += scanner.Text()
	}

	return text
}

// jsonifies the given data
func jsonify(data any) []byte {
	byteArray, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	return (byteArray)
}

var _UUID string

// called by the init function, it sets the _UUID variable 
func setupUUID() {
	uuidFilepath := filepath.Join(GetIrisDir(), "uuid")
	if CheckPathExists(uuidFilepath) {
		_UUID = readFile(uuidFilepath)
	} else {
		_UUID = uuid.New().String()
		uuidFile, err := os.Create(uuidFilepath)
		if err != nil {
			fmt.Println("unable to create uuid file")
			os.Exit(1)
		}
		if _, err = uuidFile.WriteString(_UUID); err != nil {
			fmt.Println("unable to write uuid")
			os.Exit(1)
		}
	}
}
