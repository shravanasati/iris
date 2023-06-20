package internal

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// breakIntoFrames uses ffmpeg to break the video into frames, and returns the location of
// the directory where the frames are stored, and error.
func breakIntoFrames(videoPath string) (string, error) {
	cache := loadCache()
	framesLocation, er := cache.get(videoPath)

	if er != nil {
		dirName := time.Now().Format("02-01-2006 15-04-05")
		framesLocation = filepath.Join(GetIrisDir(), "cache", dirName)
	} else {
		// video frames found from cache
		return framesLocation, er
	}

	if !CheckFileExists(framesLocation) {
		os.Mkdir(framesLocation, os.ModePerm)
	}

	cmd := exec.Command(
		"ffmpeg",
		"-i", videoPath,
		"-hide_banner",
		"-r", "5",
		"-loglevel", "debug",
		framesLocation+"/thumb_%04d.png",
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if e := cmd.Run(); e != nil {
		fmt.Println(e)
		fmt.Println(cmd.Stdout, cmd.Stderr)
		return "", fmt.Errorf("unable to break video into frames. make sure you've ffmpeg installed")
	}

	if er := cache.set(videoPath, framesLocation); er != nil {
		panic(er)
	}
	return framesLocation, nil
}

// SetVideoWallpaper sets a video as wallpaper by breaking the video into frames and then
// changes the wallpaper every few milliseconds to imitate that wallpaper is a video.
func SetVideoWallpaper(videoPath string) error {
	if !CheckFileExists(videoPath) {
		return fmt.Errorf("The file `%s` is non-existent.", videoPath)
	}
	splitted := strings.Split(videoPath, ".")
	ext := splitted[len(splitted)-1]
	if !StringInSlice(ext, []string{"mp4", "mkv"}) {
		return fmt.Errorf("The file `%s` is either unsupported or not a valid video file.", videoPath)
	}

	framesLocation, er := breakIntoFrames(videoPath)
	if er != nil {
		return er
	}

	tempConfig := &Configuration{
		WallpaperDirectory: framesLocation,
	}

	wallpapers := tempConfig.getValidWallpapers()
	sort.Strings(wallpapers)

	for {
		for _, file := range wallpapers {
			SetWallpaper(file)
			time.Sleep(time.Millisecond * 10)
		}
	}

}
