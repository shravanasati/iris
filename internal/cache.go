package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// cache repesents a cache file which keeps record of the video path and their frames directory.
type cache struct {
	location string
	data     map[string]string
}

func (c *cache) get(videoPath string) (string, error) {
	value, ok := c.data[videoPath]
	if !ok || !CheckFileExists(value) {
		delete(c.data, value) // in case checkfilexists reports false
		return "", fmt.Errorf("didn't find %v in cache", videoPath)
	}
	return value, nil
}

func (c *cache) set(videoPath, framesLocation string) error {
	c.data[videoPath] = framesLocation

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
		data:     map[string]string{},
	}

	if !CheckFileExists(cacheLocation) {
		return cacheObj
	}

	cacheContent := readFile(cacheLocation)
	if e := json.Unmarshal([]byte(cacheContent), &cacheObj.data); e != nil {
		return cacheObj
	}

	return cacheObj
}
