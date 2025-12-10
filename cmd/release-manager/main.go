package main

import (
	"flag"
	"fmt"
)

func main() {
	version := flag.String("version", "", "Release version (e.g. v1.0.0)")
	flag.Parse()
	fmt.Printf("release-manager scaffold. version=%s\n", *version)
}
