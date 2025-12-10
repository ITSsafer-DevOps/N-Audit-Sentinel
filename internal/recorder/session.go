// N-Audit Sentinel - Session Recorder
// Developer: Kristián Kašník
// Company: Nethemba s.r.o. (https://www.nethemba.com)
// License: MIT License
// Copyright (c) 2025 Kristián Kašník, Nethemba s.r.o.
// Manages PTY session with auto-respawn safety loop.
package recorder

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/term"
)

// StartSession spawns a PTY-backed shell and maintains an infinite safety loop.
// If the shell exits (e.g., user types 'exit' or Ctrl+D), the wrapper respawns it.
// The loop only exits when ctx is cancelled (e.g., via SIGUSR1).
//
// Parameters:
// - ctx: cancellation context to end the session loop.
// - logOutput: writer receiving sanitized, timestamped user input.
// - command: shell executable path (e.g., "/bin/bash").
// - args: optional arguments for the shell.
//
// Returns:
// - error: nil on normal cancellation; any shell execution error is logged and returned.
func StartSession(ctx context.Context, logOutput io.Writer, command string, args ...string) error {
	for {
		select {
		case <-ctx.Done():
			// Context cancelled, exit gracefully.
			return nil
		default:
		}

		if err := runShell(ctx, logOutput, command, args...); err != nil {
			// Log error but continue loop (respawn).
			fmt.Fprintf(os.Stderr, "[N-Audit] Shell exited with error: %v\n", err)
		}

		// Safety message after shell exit.
		fmt.Fprintln(os.Stdout, "\n[N-Audit] Session is protected. Use 'n-audit exit' to seal and close.")
	}
}

// runShell spawns a single PTY shell instance, proxies I/O, and waits for it to exit.
// It configures raw mode for stdin, handles window resize events, and ensures cleanup.
func runShell(ctx context.Context, logOutput io.Writer, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()

	// Start PTY.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("pty start failed: %w", err)
	}
	defer ptmx.Close()

	// Set initial terminal size.
	if term.IsTerminal(int(os.Stdin.Fd())) {
		if width, height, err := term.GetSize(int(os.Stdin.Fd())); err == nil {
			_ = pty.Setsize(ptmx, &pty.Winsize{Rows: uint16(height), Cols: uint16(width)})
		}
	}

	// Handle window resize signals.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	go func() {
		for range sigCh {
			if term.IsTerminal(int(os.Stdin.Fd())) {
				if width, height, err := term.GetSize(int(os.Stdin.Fd())); err == nil {
					_ = pty.Setsize(ptmx, &pty.Winsize{Rows: uint16(height), Cols: uint16(width)})
				}
			}
		}
	}()
	defer signal.Stop(sigCh)

	// Put stdin into raw mode to pass through control chars.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err == nil {
		defer term.Restore(int(os.Stdin.Fd()), oldState)
	}

	// Proxy I/O: stdin -> PTY+log, PTY -> stdout.
	done := make(chan error, 2)

	// Stdin -> PTY and log.
	go func() {
		mw := io.MultiWriter(ptmx, logOutput)
		_, err := io.Copy(mw, os.Stdin)
		done <- err
	}()

	// PTY -> Stdout.
	go func() {
		_, err := io.Copy(os.Stdout, ptmx)
		done <- err
	}()

	// Wait for command to finish or context cancel.
	cmdDone := make(chan error, 1)
	go func() {
		cmdDone <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// Kill the shell if context cancelled.
		_ = cmd.Process.Kill()
		return ctx.Err()
	case err := <-cmdDone:
		// Shell exited normally.
		return err
	}
}
