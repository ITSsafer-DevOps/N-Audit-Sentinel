# Deployment Manifest Verification Report

**Date:** 11 December 2025

## Overview

This report contains the results of attempting to validate deployment manifests using `kubectl --dry-run=client` from the current environment.

## Manifest Validation Results

Note: `kubectl` is available in the environment, but the current machine does not have a reachable Kubernetes API server (connection refused to localhost), therefore client-side OpenAPI schema validation could not be completed.

We executed a parse-only dry-run (`kubectl apply --dry-run=client --validate=false`) to check manifests; however, `kubectl` attempted to contact the API server and returned connection errors in this environment.

### Collected outputs (excerpt)

```
=== Validating (parse-only): beta-test-deployment/pod-fixed.yaml ===
E1211 07:35:12.799803  174750 memcache.go:265] couldn't get current server API group list: Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused
error: unable to recognize "beta-test-deployment/pod-fixed.yaml": Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused

=== Validating (parse-only): beta-test-deployment/serviceaccount.yaml ===
E1211 07:35:12.960334  174764 memcache.go:265] couldn't get current server API group list: Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused
unable to recognize "beta-test-deployment/serviceaccount.yaml": Get "http://localhost:8080/api?timeout=32s": dial tcp [::1]:8080: connect: connection refused
```

### File list processed

The following manifest files were checked for presence:

```
beta-test-deployment/pod-fixed.yaml
beta-test-deployment/serviceaccount.yaml
```

## Conclusion & Recommendations

- In this environment `kubectl` could not reach a Kubernetes API server, so full schema validation could not be completed. This is expected when running locally without a configured cluster or KUBECONFIG.
- Recommendation: Run the following in CI (where a kubeconfig or Kubernetes API is available) or locally with a valid `KUBECONFIG`:

```bash
for manifest in deploy/*/*.{yaml,yml} beta-test-deployment/*.{yaml,yml}; do
  echo "Validating: $manifest"
  kubectl apply --dry-run=client -f "$manifest" --validate=true | head -n 5
done
```

- Alternatively, use manifest linters such as `kubeval` or `kubeconform` in CI to validate YAML and API compatibility without accessing a live cluster.

## Status

- Result: PARTIAL — manifests present, client-side schema validation not completed due to missing API server in the current environment.
