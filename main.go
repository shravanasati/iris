package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/shravanasati/iris/cmd"
	"github.com/shravanasati/iris/internal"
)

const (
	NAME    string = "iris"
	VERSION string = "v0.4.0"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("The app has got a fatal error, and it cannot proceed further. \nPlease file a bug report at https://github.com/shravanasati/iris/issues/new/choose, with the following error message. \n```\n%s\n```", err)
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
