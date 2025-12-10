package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/backupmgr"
)

func main() {
	out := flag.String("out", "gold-master.tar.gz", "Output archive path")
	ref := flag.String("ref", "HEAD", "Git ref to archive")
	flag.Parse()
	if *out == "" {
		fmt.Fprintln(os.Stderr, "out is required")
		os.Exit(2)
	}
	if err := backupmgr.CreateSourceArchive(*ref, *out); err != nil {
		fmt.Fprintln(os.Stderr, "archive failed:", err)
		os.Exit(2)
	}
	if _, err := backupmgr.WriteChecksum(*out); err != nil {
		fmt.Fprintln(os.Stderr, "checksum failed:", err)
		os.Exit(2)
	}
	fmt.Println("Created:", *out)
}
