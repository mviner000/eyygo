package main

import (
	"fmt"
	"os"

	"github.com/mviner000/eyygo/lib"
)

func main() {
	if err := lib.ExecuteRootCommand(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
