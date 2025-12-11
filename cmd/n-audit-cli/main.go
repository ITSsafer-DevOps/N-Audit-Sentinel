// N-Audit Sentinel - Exit CLI
// Developer: Kristian Kasnik
// Company: ITSsafer-DevOps
// License: MIT License
// Copyright (c) 2025 Kristian Kasnik, ITSsafer-DevOps
// Sends SIGUSR1 to PID 1 to trigger forensic session seal and exit.
package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	// Target PID 1 (n-audit-sentinel running as container init)
	const targetPID = 1

	// Find the process
	proc, err := os.FindProcess(targetPID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[n-audit] Error: Failed to find process with PID %d: %v\n", targetPID, err)
		os.Exit(1)
	}

	// Send SIGUSR1 signal to trigger seal and exit
	if err := proc.Signal(syscall.SIGUSR1); err != nil {
		fmt.Fprintf(os.Stderr, "[n-audit] Error: Failed to send signal to PID %d: %v\n", targetPID, err)
		fmt.Fprintf(os.Stderr, "[n-audit] Hint: Ensure you have permission and the process is running.\n")
		os.Exit(1)
	}

	// Success message
	fmt.Fprintln(os.Stdout, "[n-audit] Seal signal (SIGUSR1) sent to PID 1. The session will now terminate.")
}
