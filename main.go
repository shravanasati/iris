package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/Shravan-1908/iris/cmd"
	"github.com/Shravan-1908/iris/internal"
)

const (
	NAME    string = "iris"
	VERSION string = "v0.3.0"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("The app has got a fatal error, and it cannot proceed further. \nPlease file a bug report at https://github.com/Shravan-1908/iris/issues/new/choose, with the following error message. \n```\n%s\n```", err)
			os.Exit(1)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		internal.CheckForUpdates(VERSION)
		wg.Done()
	}()

	fmt.Println(NAME, VERSION)

	cmd.Execute()
	wg.Wait()
}
