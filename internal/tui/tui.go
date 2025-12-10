// N-Audit Sentinel - TUI
// Developer: Kristián Kašník
// Company: Nethemba s.r.o. (https://www.nethemba.com)
// License: MIT License
// Copyright (c) 2025 Kristián Kašník, Nethemba s.r.o.
// Provides banner display and interactive input collection with double-enter scope logic.
package tui

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// asciiBanner is the static ASCII art banner.
var asciiBanner = `███╗   ██╗       █████╗ ██╗   ██╗██████╗ ██╗████████╗    ` + "\x1b[35m" + `███████╗███████╗███╗   ██╗████████╗██╗███╗   ██╗███████╗██╗` + "\x1b[0m" + `     
████╗  ██║      ██╔══██╗██║   ██║██╔══██╗██║╚══██╔══╝    ` + "\x1b[35m" + `██╔════╝██╔════╝████╗  ██║╚══██╔══╝██║████╗  ██║██╔════╝██║` + "\x1b[0m" + `     
██╔██╗ ██║█████╗███████║██║   ██║██║  ██║██║   ██║       ` + "\x1b[35m" + `███████╗█████╗  ██╔██╗ ██║   ██║   ██║██╔██╗ ██║█████╗  ██║` + "\x1b[0m" + `     
██║╚██╗██║╚════╝██╔══██║██║   ██║██║  ██║██║   ██║       ` + "\x1b[35m" + `╚════██║██╔══╝  ██║╚██╗██║   ██║   ██║██║╚██╗██║██╔══╝  ██║` + "\x1b[0m" + `     
██║ ╚████║      ██║  ██║╚██████╔╝██████╔╝██║   ██║       ` + "\x1b[35m" + `███████║███████╗██║ ╚████║   ██║   ██║██║ ╚████║███████╗███████╗` + "\x1b[0m" + `
╚═╝  ╚═══╝      ╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚═╝   ╚═╝       ` + "\x1b[35m" + `╚══════╝╚══════╝╚═╝  ╚═══╝   ╚═╝   ╚═╝╚═╝  ╚═══╝╚══════╝╚══════╝` + "\x1b[0m" + `
`

// ShowBanner writes the ASCII banner and project metadata to the provided writer.
// It does not perform any terminal control and is safe for non-TTY outputs.
func ShowBanner(w io.Writer) {
	fmt.Fprintln(w, asciiBanner)
	fmt.Fprintln(w, "Developer: Kristián Kašník")
	fmt.Fprintln(w, "Company: ITSsafer-DevOps")
	fmt.Fprintln(w, "License: MIT License (Open Source)")
	fmt.Fprintln(w)
}

// GetPentesterInfo prompts for pentester and client names, trimming whitespace.
// It reads two lines from r and writes prompts to w.
// Returns the collected names or io.EOF if input stream ends.
func GetPentesterInfo(r io.Reader, w io.Writer) (pentesterName, clientName string, err error) {
	scanner := bufio.NewScanner(r)
	fmt.Fprint(w, "Pentester Name: ")
	if !scanner.Scan() {
		if scanner.Err() != nil {
			return "", "", scanner.Err()
		}
		return "", "", io.EOF
	}
	pentesterName = strings.TrimSpace(scanner.Text())
	fmt.Fprint(w, "Client Name: ")
	if !scanner.Scan() {
		if scanner.Err() != nil {
			return "", "", scanner.Err()
		}
		return pentesterName, "", io.EOF
	}
	clientName = strings.TrimSpace(scanner.Text())
	return pentesterName, clientName, nil
}

// GetScope collects IP/CIDR and domain scope using double-enter transitions.
// Two consecutive empty lines end the current loop.
// First loop collects IP/CIDR; second loop collects Domains.
// Note: raw inputs are returned; validation occurs in the orchestrator.
func GetScope(r io.Reader, w io.Writer) (ips, domains []string, err error) {
	scanner := bufio.NewScanner(r)

	// Layer 3 (IP/CIDR) loop
	fmt.Fprintln(w, "[Layer 3] Enter Target IP/CIDR (Double Enter to continue):")
	prevEmpty := false
	for {
		fmt.Fprint(w, "IP/CIDR> ")
		if !scanner.Scan() {
			if scanner.Err() != nil {
				return ips, domains, scanner.Err()
			}
			return ips, domains, nil
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if prevEmpty { // double enter -> exit loop
				break
			}
			prevEmpty = true
			continue
		}
		prevEmpty = false
		ips = append(ips, line)
	}

	// Layer 7 (Domain) loop
	fmt.Fprintln(w, "[Layer 7] Enter Target Domain (Double Enter to finish):")
	prevEmpty = false
	for {
		fmt.Fprint(w, "Domain> ")
		if !scanner.Scan() {
			if scanner.Err() != nil {
				return ips, domains, scanner.Err()
			}
			return ips, domains, nil
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if prevEmpty { // double enter -> end function
				break
			}
			prevEmpty = true
			continue
		}
		prevEmpty = false
		domains = append(domains, line)
	}
	return ips, domains, nil
}
