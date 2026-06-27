package main

import (
	"fmt"
	"os"
	"main.go/helpers"
	"main.go/constants"
)

func main() {
	c := make(chan helpers.Result)

	for i, url := range constants.URLs {
		go helpers.CheckURL(url, i, c)
	}

	results := make([]string, len(constants.URLs))
	for i := 0; i < len(constants.URLs); i++ {
		r := <-c
		results[r.Index] = r.Msg
	}

	for _, msg := range results {
		fmt.Println(msg)
	}

	if os.Getenv("DEV") == "true" {
		fmt.Println("(ran in DEV mode)")
	}
}