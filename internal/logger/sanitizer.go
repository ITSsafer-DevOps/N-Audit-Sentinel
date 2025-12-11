// N-Audit Sentinel - Logger Sanitizer
// Developer: Kristian Kasnik
// Company: ITSsafer-DevOps
// License: MIT License
// Copyright (c) 2025 Kristian Kasnik, ITSsafer-DevOps
// This file provides ANSI escape sequence stripping and a timestamping writer.
package logger

import (
	"bytes"
	"io"
	"regexp"
	"sync"
	"time"
)

// Precompiled regex matching ANSI escape sequences (CSI, OSC, single chars).
// Pattern rationale:
// \x1b followed by either:
//   - '[' then zero or more parameter bytes (0-? and space-/) then a final byte @-~
//   - a single 7-bit C1 control in range @-Z\\-_
//
// This covers common color codes, cursor movement, erase commands, etc.
var ansiRegex = regexp.MustCompile("\x1b(?:\\[[0-?]*[ -/]*[@-~]|[@-Z\\-_])")

// StripANSI removes ANSI escape sequences from the provided byte slice.
// It is intended for sanitizing PTY output prior to timestamping and logging.
func StripANSI(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	cleaned := ansiRegex.ReplaceAll(data, []byte{})
	return cleaned
}

// TimestampedWriter wraps an io.Writer and prepends RFC3339Nano timestamps
// (UTC) to each complete line written.
//
// Behavior:
// - Partial lines are buffered until a newline is received.
// - Each line is sanitized via StripANSI and written with a timestamp prefix.
// - Write is concurrency-safe via internal mutex.
type TimestampedWriter struct {
	w   io.Writer
	buf bytes.Buffer
	mu  sync.Mutex
}

// NewTimestampedWriter constructs a new TimestampedWriter.
// w: the underlying destination writer to receive timestamped, sanitized lines.
func NewTimestampedWriter(w io.Writer) *TimestampedWriter {
	return &TimestampedWriter{w: w}
}

// Write implements io.Writer.
// It buffers incoming data until newline boundaries are found.
// Each complete line is stripped of ANSI sequences and written with
// a timestamp prefix: "<RFC3339Nano> <line>".
func (tw *TimestampedWriter) Write(p []byte) (n int, err error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	// Append incoming data to buffer.
	if _, err = tw.buf.Write(p); err != nil {
		return 0, err
	}

	for {
		data := tw.buf.Bytes()
		idx := bytes.IndexByte(data, '\n')
		if idx == -1 {
			break // no complete line yet
		}
		line := data[:idx+1] // include newline
		// Advance buffer (consume processed line)
		tw.buf.Next(idx + 1)

		// Sanitize and timestamp line.
		clean := StripANSI(line)
		ts := time.Now().UTC().Format(time.RFC3339Nano) + " "
		if _, err = tw.w.Write([]byte(ts)); err != nil {
			return len(p), err
		}
		if _, err = tw.w.Write(clean); err != nil {
			return len(p), err
		}
	}
	return len(p), nil
}
