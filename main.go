package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	tableMode := flag.Bool("table", false, "show output in table format")
	flag.BoolVar(tableMode, "t", false, "show output in table format (shorthand)")
	noColor := flag.Bool("no-color", false, "disable colored output")
	flag.Parse()

	initColor(*noColor)

	token, err := resolveToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	result := fetchAll(token)

	if *tableMode {
		displayTable(result)
	} else {
		displayCompact(result)
	}

	if result.MyPRsErr != nil && result.ReviewReqErr != nil {
		os.Exit(1)
	}
}
