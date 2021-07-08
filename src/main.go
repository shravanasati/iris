package main

import (
	"github.com/reujab/wallpaper"
)

func main() {
	// err = wallpaper.SetFromFile(`C:\Users\LENOVO\Downloads\Images\wp9310707-summer-scotland-wallpapers.jpg`)
	// if err != nil {
	// 	panic(err)
	// }

	err := wallpaper.SetFromURL("https://source.unsplash.com/1600x900/?nature")
	if err != nil {
		panic(err)
	}

}