This `pkg/` directory is reserved for libraries intended for external use.

Project layout guidance:

- `cmd/` contains main applications for this project.
- `internal/` contains packages only for use inside this repository.
- `pkg/` is optional and should contain code intended to be imported by external projects.

This file was added as part of a minimal alignment with Google Go project conventions during the final cleanup and audit.
