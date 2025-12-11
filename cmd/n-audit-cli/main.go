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

// processSignaler abstracts an entity that can receive an OS signal. `*os.Process`
// satisfies this interface which makes the helper testable by injecting fakes.
type processSignaler interface {
	Signal(os.Signal) error
}

// finder returns a `processSignaler` for a PID, or an error.
type finder func(pid int) (processSignaler, error)

// SendSealSignalWithFinder sends SIGUSR1 to the given PID using the provided
// `finder`. It returns an error on lookup or signal failure.
func SendSealSignalWithFinder(find finder, pid int) error {
	proc, err := find(pid)
	if err != nil {
		return fmt.Errorf("failed to find process %d: %w", pid, err)
	}
	if err := proc.Signal(syscall.SIGUSR1); err != nil {
		return fmt.Errorf("failed to send signal to %d: %w", pid, err)
	}
	return nil
}

func main() {
	// Target PID 1 (n-audit-sentinel running as container init)
	const targetPID = 1

	// Use os.FindProcess as the production finder. os.FindProcess returns
	// *os.Process which implements processSignaler.
	osFinder := func(pid int) (processSignaler, error) { return os.FindProcess(pid) }

	if err := SendSealSignalWithFinder(osFinder, targetPID); err != nil {
		fmt.Fprintf(os.Stderr, "[n-audit] Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "[n-audit] Hint: Ensure you have permission and the process is running.\n")
		os.Exit(1)
	}

	// Success message
	fmt.Fprintln(os.Stdout, "[n-audit] Seal signal (SIGUSR1) sent to PID 1. The session will now terminate.")
}
