package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
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
	LogInfof("cache", "retrieving cache for: %s", videoPath)
	value, ok := c.data[videoPath]
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		LogErrorf("cache", "failed to stat video file %s: %v", videoPath, err)
		return "", err
	}
	fileLastModTime := fileInfo.ModTime().Format(timeFormat)
	if !ok {
		LogInfof("cache", "video not found in cache: %s", videoPath)
		return "", fmt.Errorf("didn't find %v in cache", videoPath)
	}
	if fileLastModTime != value.LastMtime {
		LogWarnf("cache", "cache expired for %s (mtime mismatch)", videoPath)
		return "", fmt.Errorf("%v expired", videoPath)
	}
	if !CheckPathExists(value.FramesFolderPath) {
		LogWarnf("cache", "frames folder missing for %s: %s", videoPath, value.FramesFolderPath)
		delete(c.data, videoPath) // in case checkfilexists reports false
		return "", fmt.Errorf("didn't find %v in cache", videoPath)
	}
	LogInfof("cache", "cache hit for %s at %s", videoPath, value.FramesFolderPath)
	return value.FramesFolderPath, nil
}

func (c *cache) set(videoPath, framesLocation string) error {
	LogInfof("cache", "caching frames for %s at %s", videoPath, framesLocation)
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		LogErrorf("cache", "failed to stat video file %s: %v", videoPath, err)
		return err
	}
	c.data[videoPath] = cacheEntry{
		FramesFolderPath: framesLocation,
		LastMtime:        fileInfo.ModTime().Format(timeFormat),
	}

	return c.write()
}

func (c *cache) write() error {
	LogInfof("cache", "writing cache to: %s", c.location)
	f, err := os.Create(c.location)
	if err != nil {
		LogErrorf("cache", "failed to create cache file: %v", err)
		return err
	}
	defer f.Close()
	_, wr := f.Write(jsonify(c.data))
	if wr != nil {
		LogErrorf("cache", "failed to write cache data: %v", wr)
	}
	return wr
}

func loadCache() *cache {
	cacheLocation := filepath.Join(GetIrisDir(), "cache", "cache.json")
	LogInfof("cache", "loading cache from: %s", cacheLocation)
	cacheObj := &cache{
		location: cacheLocation,
		data:     cacheEntryMap{},
	}

	if !CheckPathExists(cacheLocation) {
		LogInfof("cache", "no cache file found")
		return cacheObj
	}

	cacheContent := readFile(cacheLocation)
	if e := json.Unmarshal([]byte(cacheContent), &cacheObj.data); e != nil {
		LogErrorf("cache", "failed to unmarshal cache: %v", e)
		return cacheObj
	}

	LogInfof("cache", "cache loaded with %d entries", len(cacheObj.data))
	return cacheObj
}

// CacheSize returns total cache size, including video frames and remote source results.
func CacheSize() ByteSize {
	cacheLocation := filepath.Join(GetIrisDir(), "cache")
	var size int64

	err := filepath.Walk(cacheLocation, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0
	}
	return ByteSize(size)
}

// CacheClear empties all iris cache, including videos and remote source results.
func CacheClear() error {
	LogInfof("cache", "clearing all cache")
	cacheLocation := filepath.Join(GetIrisDir(), "cache")
	subdirs, err := filepath.Glob(filepath.Join(cacheLocation, "*"))
	if err != nil {
		LogErrorf("cache", "failed to glob cache directory: %v", err)
		return err
	}

	// Iterate over subdirectories and remove them
	for _, subdir := range subdirs {
		fileInfo, err := os.Stat(subdir)
		if err != nil {
			continue
		}

		if fileInfo.IsDir() {
			LogInfof("cache", "removing directory: %s", subdir)
			// If it's a directory, remove it
			err := os.RemoveAll(subdir)
			if err != nil {
				LogErrorf("cache", "failed to remove directory %s: %v", subdir, err)
				return err
			}
		} else if fileInfo.Name() == "github.json" {
			LogInfof("cache", "removing github cache file")
			// Also remove github cache
			if err := os.Remove(subdir); err != nil {
				LogErrorf("cache", "failed to remove github cache: %v", err)
				return err
			}
		}
	}

	// remove all references from cache.json
	ca := &cache{
		location: filepath.Join(cacheLocation, "cache.json"),
		data:     cacheEntryMap{},
	}
	if err = ca.write(); err != nil {
		LogErrorf("cache", "failed to clear cache.json references: %v", err)
		return err
	}

	LogInfof("cache", "cache cleared successfully")
	return nil
}

// CacheRemove removes a single item from the iris video cache.
func CacheRemove(videoPath string) error {
	LogInfof("cache", "removing item from cache: %s", videoPath)
	ca := loadCache()
	val, ok := ca.data[videoPath]
	if !ok {
		LogWarnf("cache", "video not found in cache for removal: %s", videoPath)
		return fmt.Errorf("video %s not found in cache", videoPath)
	}

	if CheckPathExists(val.FramesFolderPath) {
		LogInfof("cache", "removing frames folder: %s", val.FramesFolderPath)
		err := os.RemoveAll(val.FramesFolderPath)
		if err != nil {
			LogErrorf("cache", "failed to remove frames folder %s: %v", val.FramesFolderPath, err)
			return err
		}
	}

	delete(ca.data, videoPath)
	return ca.write()
}

func CacheShow() {
	ca := loadCache()
	tableData := [][]string{}
	i := 1
	for videoPath := range ca.data {
		tableData = append(tableData, []string{fmt.Sprintf("%v", i), videoPath})
		i++
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"S. No.", "Video"})
	table.AppendBulk(tableData)
	table.Render()
}
