// N-Audit Sentinel - Kubernetes Discovery
// Developer: Kristián Kašník
// Company: Nethemba s.r.o. (https://www.nethemba.com)
// License: MIT License
// Copyright (c) 2025 Kristián Kašník, Nethemba s.r.o.
// Automatic discovery of Kubernetes API server and DNS servers.
package discovery

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// DiscoverK8sAPI discovers the Kubernetes API server endpoint from environment variables.
// Returns the combined endpoint as "host:port" or an error if not found.
func DiscoverK8sAPI() (string, error) {
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")

	if host == "" || port == "" {
		return "", fmt.Errorf("KUBERNETES_SERVICE_HOST or KUBERNETES_SERVICE_PORT not set")
	}

	return fmt.Sprintf("%s:%s", host, port), nil
}

// DiscoverDNS parses the resolv.conf file and extracts nameserver IP addresses.
// Returns a slice of DNS server IPs.
func DiscoverDNS(resolvConfPath string) ([]string, error) {
	file, err := os.Open(resolvConfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", resolvConfPath, err)
	}
	defer file.Close()

	var nameservers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "nameserver") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				nameservers = append(nameservers, fields[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading %s: %w", resolvConfPath, err)
	}

	return nameservers, nil
}
