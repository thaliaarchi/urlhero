package main

import (
	"fmt"
	"os"

	"github.com/andrewarchi/urlteam/wwiki"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s dir\n", os.Args[0])
		os.Exit(2)
	}
	err := wwiki.DownloadDumps(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
