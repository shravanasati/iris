package main

import (
	"fmt"
	"os"

	"github.com/Shravan-1908/iris/cmd"
)

const (
	NAME    string = "iris"
	VERSION string = "v0.2.1"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("The app has got a fatal error, and it cannot proceed further. \nPlease file a bug report at https://github.com/Shravan-1908/issues/new/choose, with the following error message. \n```\n%s\n```", err)
			os.Exit(1)
		}
	}()

	fmt.Println(NAME, VERSION)

	cmd.Execute()
}
