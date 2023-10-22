package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type cacheEntry struct {
	LastMtime        string `json:"last_modified_time"` // video's last modified time, used to compare and refresh the cache if needed
	FramesFolderPath string `json:"frames_folder_path"` // the absolute path to the frames folder of the video
}

type cacheEntryMap map[string]cacheEntry

// cache repesents a cache file which keeps record of the video path and their frames directory.
type cache struct {
	location string
	data     cacheEntryMap // a map of video paths and resp cache entries
}

func (c *cache) get(videoPath string) (string, error) {
	value, ok := c.data[videoPath]
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return "", err
	}
	fileLastModTime := fileInfo.ModTime().Format(timeFormat)
	if fileLastModTime != value.LastMtime {
		return "", fmt.Errorf("%v expired", videoPath)
	}
	if !ok || !CheckPathExists(value.FramesFolderPath) {
		delete(c.data, videoPath) // in case checkfilexists reports false
		return "", fmt.Errorf("didn't find %v in cache", videoPath)
	}
	// todo add logic to check for LastMtime
	return value.FramesFolderPath, nil
}

func (c *cache) set(videoPath, framesLocation string) error {
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return err
	}
	c.data[videoPath] = cacheEntry{
		FramesFolderPath: framesLocation,
		LastMtime:        fileInfo.ModTime().Format(timeFormat),
	}

	f, err := os.Create(c.location)
	if err != nil {
		return err
	}
	defer f.Close()
	_, wr := f.Write(jsonify(c.data))
	return wr
}

func loadCache() *cache {
	cacheLocation := filepath.Join(GetIrisDir(), "cache", "cache.json")
	cacheObj := &cache{
		location: cacheLocation,
		data:     cacheEntryMap{},
	}

	if !CheckPathExists(cacheLocation) {
		return cacheObj
	}

	cacheContent := readFile(cacheLocation)
	if e := json.Unmarshal([]byte(cacheContent), &cacheObj.data); e != nil {
		return cacheObj
	}

	return cacheObj
}

// Returns total cache size.
func CacheSize() ByteSize {
	cacheLocation := filepath.Join(GetIrisDir(), "cache")
	var size int64

	err := filepath.Walk(cacheLocation, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) != "json" {
			size += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0
	}
	return ByteSize(size)
}

// CacheEmpty empties all iris video cache.
func CacheEmpty() error {
	cacheLocation := filepath.Join(GetIrisDir(), "cache")
	subdirs, err := filepath.Glob(filepath.Join(cacheLocation, "*"))
	if err != nil {
		return err
	}

	// Iterate over subdirectories and remove them
	for _, subdir := range subdirs {
		fileInfo, err := os.Stat(subdir)
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			// If it's a directory, remove it
			err := os.RemoveAll(subdir)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
