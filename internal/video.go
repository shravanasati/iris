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

var AllowedVideoExtensions = []string{"mp4", "mkv", "gif"}

// breakIntoFrames uses ffmpeg to break the video into frames, and returns the location of
// the directory where the frames are stored, and error.
func breakIntoFrames(videoPath string) (string, error) {
	LogInfof("video", "preparing to break video into frames: %s", videoPath)
	cache := loadCache()
	framesLocation, er := cache.get(videoPath)

	if er != nil {
		LogInfof("video", "cache miss for video, generating new frames directory")
		dirName := time.Now().Format("02-01-2006 15-04-05")
		framesLocation = filepath.Join(GetIrisDir(), "cache", dirName)
	} else {
		LogInfof("video", "using cached frames from: %s", framesLocation)
		return framesLocation, er
	}

	LogInfof("video", "breaking video into frames for the first time, this may take a while")

	if !CheckPathExists(framesLocation) {
		LogInfof("video", "creating frames directory: %s", framesLocation)
		os.Mkdir(framesLocation, os.ModePerm)
	}

	LogInfof("video", "running ffmpeg for: %s", videoPath)
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
		LogErrorf("video", "ffmpeg failed for %s: %v", videoPath, e)
		LogErrorf("video", "ffmpeg error output: %s", stderr.String())
		fmt.Println(e)
		fmt.Println(cmd.Stdout, cmd.Stderr)
		return "", fmt.Errorf("unable to break video into frames. make sure you've ffmpeg installed and present on path")
	}

	LogInfof("video", "ffmpeg completed successfully")
	absPath, err := filepath.Abs(videoPath)
	if err != nil {
		absPath = videoPath
	}
	if er := cache.set(absPath, framesLocation); er != nil {
		LogErrorf("video", "failed to update cache for video: %v", er)
		panic(er)
	}
	return framesLocation, nil
}

// SetVideoWallpaper sets a video as wallpaper by breaking the video into frames and then
// changes the wallpaper every few milliseconds to imitate that wallpaper is a video.
func SetVideoWallpaper(videoPath string) error {
	LogInfof("video", "setting video wallpaper: %s", videoPath)
	if !CheckPathExists(videoPath) {
		LogErrorf("video", "video file does not exist: %s", videoPath)
		return fmt.Errorf("the file `%s` is non-existent", videoPath)
	}
	splitted := strings.Split(videoPath, ".")
	ext := splitted[len(splitted)-1]
	if !ItemInSlice(ext, AllowedVideoExtensions) {
		LogErrorf("video", "unsupported video extension for %s: %s", videoPath, ext)
		return fmt.Errorf("the file `%s` is either unsupported or not a valid video file", videoPath)
	}

	framesLocation, er := breakIntoFrames(videoPath)
	if er != nil {
		LogErrorf("video", "failed to break video into frames: %v", er)
		return er
	}

	tempConfig := &Configuration{
		WallpaperDirectory: framesLocation,
	}

	LogInfof("video", "loading wallpapers from frames directory")
	wallpapers := tempConfig.getValidWallpapers()
	sort.Strings(wallpapers)
	LogInfof("video", "playing video wallpaper with %d frames", len(wallpapers))

	for {
		for _, file := range wallpapers {
			SetWallpaper(file)
			time.Sleep(time.Millisecond * 10)
		}
	}

}
