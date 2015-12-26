package main

import (
	"flag"
	"os"
	"fmt"
	"github.com/aoeu/brainguck"
)

func main() {
	var filename string
	var verbose bool
	flag.StringVar(&filename, "in", "", "The source code filename to use as input.")
	flag.BoolVar(&verbose,"v", false, "Verbosely print the number of bytes of input interpreted.")
	flag.Parse()
	if filename == "" {
		os.Exit(1)
	}
	n, err := brainguck.InterpretFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if verbose {
		fmt.Printf("%v bytes read.\n", n)
	}
}
