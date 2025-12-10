# Multi-stage build for N-Audit Sentinel

# --- Build Stage ---
FROM golang:1.24-alpine AS build
WORKDIR /src

# Enable modules
ENV CGO_ENABLED=0

# Placeholder go.mod (will be refined later if not present)
# Copy module files first for caching
COPY go.mod go.sum* /src/
RUN go mod download || true

# Copy source
COPY . /src

# Build binaries
RUN go build -o /out/n-audit-sentinel ./cmd/n-audit-sentinel \
 && go build -o /out/n-audit ./cmd/n-audit-cli

# --- Final Stage ---
FROM kalilinux/kali-rolling AS final

# Create non-root user/group 'auditor'
RUN groupadd -r auditor && useradd -r -g auditor -m auditor

# Install minimal runtime dependencies (optional - kept lean)
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates bash \
    && rm -rf /var/lib/apt/lists/*

# Copy binaries
COPY --from=build /out/n-audit-sentinel /usr/local/bin/n-audit-sentinel
COPY --from=build /out/n-audit /usr/local/bin/n-audit

# Ensure executables have correct permissions
RUN chmod 0755 /usr/local/bin/n-audit-sentinel /usr/local/bin/n-audit

# Default working directory for logs/data (PVC mount target suggested)
WORKDIR /var/lib/n-audit

# Ensure runtime directory exists with strict permissions and ownership
RUN mkdir -p /var/lib/n-audit \
    && chown auditor:auditor /var/lib/n-audit \
    && chmod 0700 /var/lib/n-audit

# NOTE: Runtime may override user to root via Pod securityContext if needed.
USER auditor

ENTRYPOINT ["/usr/local/bin/n-audit-sentinel"]
