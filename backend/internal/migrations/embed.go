package migrations

import "embed"

// FS holds all SQL migration files bundled into the binary.
//
//go:embed sql/*.sql
var FS embed.FS
