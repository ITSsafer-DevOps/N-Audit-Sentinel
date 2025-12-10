package main

import (
	"flag"
	"fmt"
)

func main() {
	out := flag.String("out", "", "Output archive path")
	flag.Parse()
	fmt.Printf("backup-manager scaffold. out=%s\n", *out)
}
