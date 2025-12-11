// N-Audit Sentinel
// Developer: Kristian Kasnik
// Company: ITSsafer-DevOps
// License: MIT License
// Copyright (c) 2025 Kristian Kasnik, ITSsafer-DevOps
// Main PID 1 wrapper orchestrating forensic session lifecycle.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/discovery"
	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/logger"
	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/policy"
	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/recorder"
	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/signature"
	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/tui"
	"github.com/ITSsafer-DevOps/N-Audit-Sentinel/internal/validation"
)

const (
	logDir  = "/var/lib/n-audit"
	logFile = "session.log"
)

func main() {
	// Setup signal handling and cancellable context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGUSR1, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		fmt.Fprintf(os.Stdout, "\n[N-Audit] Received signal %s. Initiating shutdown...\n", sig)
		cancel()
	}()

	// Debug mode: wait for manual attachment
	if os.Getenv("N_AUDIT_DEBUG") == "true" {
		const debugWait = 30 * time.Second
		fmt.Fprintf(os.Stdout, "[N-Audit] DEBUG MODE ACTIVE\n")
		fmt.Fprintf(os.Stdout, "[N-Audit] Attach BEFORE countdown ends to interact with TUI.\n")
		fmt.Fprintf(os.Stdout, "[N-Audit] Use: kubectl attach -it n-audit-sentinel -c n-audit\n")
		fmt.Fprintf(os.Stdout, "[N-Audit] (Avoid kubectl exec – it starts a separate process.)\n")
		deadline := time.Now().Add(debugWait)
		for {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				fmt.Fprintf(os.Stdout, "[N-Audit] Debug countdown finished. Starting TUI...\n")
				break
			}
			// Print every 5 seconds (and at start)
			if int(remaining.Seconds())%5 == 0 || remaining > debugWait-1*time.Second {
				fmt.Fprintf(os.Stdout, "[N-Audit] Countdown: %2.0fs remaining...\n", remaining.Seconds())
			}
			select {
			case <-ctx.Done():
				fmt.Fprintf(os.Stdout, "[N-Audit] Debug wait interrupted by signal.\n")
				return
			case <-time.After(1 * time.Second):
			}
		}
	}

	// Discover Kubernetes environment.
	apiServer, err := discovery.DiscoverK8sAPI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[N-Audit] Warning: Failed to discover K8s API: %v\n", err)
		apiServer = "unknown"
	} else {
		fmt.Fprintf(os.Stdout, "[N-Audit] Discovered K8s API Server: %s\n", apiServer)
	}

	dnsServers, err := discovery.DiscoverDNS("/etc/resolv.conf")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[N-Audit] Warning: Failed to discover DNS servers: %v\n", err)
		dnsServers = []string{}
	} else {
		fmt.Fprintf(os.Stdout, "[N-Audit] Discovered DNS Servers: %v\n", dnsServers)
	}

	// Initialize logging infrastructure.
	if err := os.MkdirAll(logDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "[N-Audit] Error: Failed to create log directory: %v\n", err)
		os.Exit(1)
	}

	logPath := filepath.Join(logDir, logFile)
	logFileHandle, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[N-Audit] Error: Failed to open log file: %v\n", err)
		os.Exit(1)
	}
	// Will close explicitly during teardown before sealing

	// Wrap log file with timestamped writer.
	timestampedLog := logger.NewTimestampedWriter(logFileHandle)

	// Log discovered infrastructure.
	fmt.Fprintf(timestampedLog, "\n=== Infrastructure Discovery ===\n")
	fmt.Fprintf(timestampedLog, "K8s API Server: %s\n", apiServer)
	fmt.Fprintf(timestampedLog, "DNS Servers: %v\n", dnsServers)
	fmt.Fprintf(timestampedLog, "================================\n\n")

	// Execute TUI flow: banner, identification, scope.
	tui.ShowBanner(os.Stdout)

	pentester, client, err := tui.GetPentesterInfo(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[N-Audit] Error: Failed to get pentester info: %v\n", err)
		os.Exit(1)
	}

	// Write session header to log.
	fmt.Fprintf(timestampedLog, "=== N-Audit Sentinel Session ===\n")
	fmt.Fprintf(timestampedLog, "Pentester: %s\n", pentester)
	fmt.Fprintf(timestampedLog, "Client: %s\n", client)
	fmt.Fprintf(timestampedLog, "================================\n\n")

	ips, domains, err := tui.GetScope(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[N-Audit] Error: Failed to get scope: %v\n", err)
		os.Exit(1)
	}

	// Validate and normalize scope entries via helper.
	validIPs, validDomains, warnings := validation.ValidateScope(ips, domains)
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "[N-Audit] Warning: %s\n", w)
		fmt.Fprintf(timestampedLog, "Warning: %s\n", w)
	}

	// Log validated scope.
	fmt.Fprintf(timestampedLog, "--- Scope Definition ---\n")
	if len(validIPs) > 0 {
		fmt.Fprintf(timestampedLog, "Target IPs/CIDR: %v\n", validIPs)
	}
	if len(validDomains) > 0 {
		fmt.Fprintf(timestampedLog, "Target Domains: %v\n", validDomains)
	}
	if len(ips) == 0 && len(domains) == 0 {
		fmt.Fprintf(timestampedLog, "Scope: Unrestricted Mode\n")
		fmt.Fprintln(os.Stdout, "[N-Audit] No scope defined. Running in unrestricted mode.")
	}
	fmt.Fprintf(timestampedLog, "\n")

	// Apply Cilium Network Policy only if scope provided.
	if len(validIPs) > 0 || len(validDomains) > 0 {
		// IMPORTANT: Ensure Terraform deployment adds label "app: n-audit-sentinel" to the pod.
		namespace := os.Getenv("POD_NAMESPACE")
		if namespace == "" {
			namespace = "default"
		}

		ciliumClient, err := policy.NewCiliumClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[N-Audit] Warning: Failed to create Cilium client: %v\n", err)
			fmt.Fprintf(timestampedLog, "Warning: Policy enforcement unavailable - %v\n", err)
		} else {
			podLabels := map[string]string{"app": "n-audit-sentinel"}
			policyName := "n-audit-policy"

			fmt.Fprintln(os.Stdout, "[N-Audit] Applying network policy (3-zone enforcement)...")
			fmt.Fprintf(timestampedLog, "--- Network Policy ---\n")
			fmt.Fprintf(timestampedLog, "Policy Name: %s\n", policyName)

			if err := ciliumClient.ApplyPolicy(policyName, namespace, podLabels, dnsServers, apiServer, validIPs, validDomains); err != nil {
				fmt.Fprintf(os.Stderr, "[N-Audit] Warning: Failed to apply policy: %v\n", err)
				fmt.Fprintf(timestampedLog, "Status: FAILED - %v\n\n", err)
			} else {
				fmt.Fprintln(os.Stdout, "[N-Audit] Network policy applied successfully.")
				fmt.Fprintf(timestampedLog, "Status: APPLIED\n")
				fmt.Fprintf(timestampedLog, "Zones: Infra(API+DNS) + Maintenance(*.kali.org,github.com,docker.io,gitlab.com,pypi.org,crates.io) + Target(scope)\n\n")
			}

			// Ensure policy cleanup on exit.
			defer func() {
				fmt.Fprintln(os.Stdout, "[N-Audit] Removing network policy...")
				if err := ciliumClient.DeletePolicy(policyName, namespace); err != nil {
					fmt.Fprintf(os.Stderr, "[N-Audit] Warning: Failed to delete policy: %v\n", err)
					fmt.Fprintf(timestampedLog, "Warning: Policy deletion failed - %v\n", err)
				} else {
					fmt.Fprintf(timestampedLog, "Policy removed: %s\n", policyName)
				}
			}()
		}
	}

	// Start the protected session loop.
	fmt.Fprintln(os.Stdout, "[N-Audit] Starting protected session...")
	fmt.Fprintf(timestampedLog, "--- Session Started ---\n\n")

	if err := recorder.StartSession(ctx, timestampedLog, "/bin/bash"); err != nil && err != context.Canceled {
		fmt.Fprintf(os.Stderr, "[N-Audit] Session error: %v\n", err)
	}

	// Teardown: log seal message.
	fmt.Fprintln(os.Stdout, "\n[N-Audit] Session terminated. Sealing log...")
	fmt.Fprintf(timestampedLog, "\n--- Session Ended ---\n")

	// IMPORTANT: Do not write to timestampedLog after closing the file handle.
	// Close log file to flush all content before sealing.
	if err := logFileHandle.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "[N-Audit] Warning: Failed to close log file before sealing: %v\n", err)
	}
	// Seal log if SSH private key is provided.
	keyPath := os.Getenv("SSH_SIGN_KEY_PATH")
	if keyPath == "" {
		fmt.Fprintln(os.Stdout, "[N-Audit] Warning: SSH_SIGN_KEY_PATH not set; skipping seal.")
	} else {
		if err := signature.SealLogFile(logPath, keyPath); err != nil {
			fmt.Fprintf(os.Stderr, "[N-Audit] Warning: Failed to seal log: %v\n", err)
		} else {
			fmt.Fprintln(os.Stdout, "[N-Audit] Log sealed with SSH signature.")
		}
	}

	fmt.Fprintln(os.Stdout, "[N-Audit] Log sealed. Exiting.")

	// (Policy cleanup handled via defer if applied earlier)
}
