# Runtime image for beta test using release binaries
FROM kalilinux/kali-rolling
RUN groupadd -r auditor && useradd -r -g auditor -m auditor \
 && apt-get update && apt-get install -y --no-install-recommends ca-certificates bash \
 && rm -rf /var/lib/apt/lists/*
COPY n-audit-sentinel /usr/local/bin/n-audit-sentinel
COPY n-audit /usr/local/bin/n-audit
RUN chmod 0755 /usr/local/bin/n-audit-sentinel /usr/local/bin/n-audit \
 && mkdir -p /var/lib/n-audit && chown auditor:auditor /var/lib/n-audit && chmod 0700 /var/lib/n-audit
WORKDIR /var/lib/n-audit
USER auditor
ENTRYPOINT ["/usr/local/bin/n-audit-sentinel"]
