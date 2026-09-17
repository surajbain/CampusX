package migrations

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

const migrationsTable = "schema_migrations"

// Run applies all pending migrations in lexical order.
// Each migration runs in its own transaction; on success a row is recorded
// in schema_migrations with its checksum.
//
// SQL files are split into individual statements because Postgres does not
// allow multiple commands in a single prepared statement.
func Run(ctx context.Context, gdb *gorm.DB, log *zerolog.Logger) error {
	if err := ensureMigrationsTable(ctx, gdb); err != nil {
		return fmt.Errorf("ensure migrations table: %w", err)
	}

	applied, err := loadApplied(ctx, gdb)
	if err != nil {
		return fmt.Errorf("load applied: %w", err)
	}

	files, err := listMigrationFiles()
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}

	if len(files) == 0 {
		log.Warn().Msg("no migration files found")
		return nil
	}

	pending := 0
	for _, file := range files {
		name := file
		body, err := FS.ReadFile("sql/" + name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		sum := checksum(body)

		if prev, ok := applied[name]; ok {
			if prev != sum {
				log.Warn().
					Str("migration", name).
					Str("recorded", prev).
					Str("computed", sum).
					Msg("migration checksum mismatch — file changed after apply")
			}
			continue
		}

		log.Info().Str("migration", name).Msg("applying migration")

		err = gdb.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			stmts := splitSQLStatements(string(body))
			for i, stmt := range stmts {
				if err := tx.Exec(stmt).Error; err != nil {
					return fmt.Errorf(
						"exec statement %d/%d: %w\nSQL: %s",
						i+1, len(stmts), err, truncate(stmt, 300),
					)
				}
			}
			return tx.Exec(
				"INSERT INTO "+migrationsTable+" (name, checksum, applied_at) VALUES (?, ?, ?)",
				name, sum, time.Now().UTC(),
			).Error
		})
		if err != nil {
			return fmt.Errorf("apply %s: %w", name, err)
		}

		log.Info().Str("migration", name).Msg("migration applied")
		pending++
	}

	if pending == 0 {
		log.Info().Msg("database schema already up-to-date")
	} else {
		log.Info().Int("count", pending).Msg("migrations applied")
	}
	return nil
}

// ensureMigrationsTable creates the tracking table if missing.
func ensureMigrationsTable(ctx context.Context, gdb *gorm.DB) error {
	stmt := `
	CREATE TABLE IF NOT EXISTS ` + migrationsTable + ` (
		name        VARCHAR(200) PRIMARY KEY,
		checksum    VARCHAR(64) NOT NULL,
		applied_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`
	return gdb.WithContext(ctx).Exec(stmt).Error
}

func loadApplied(ctx context.Context, gdb *gorm.DB) (map[string]string, error) {
	type row struct {
		Name     string
		Checksum string
	}
	var rows []row
	if err := gdb.WithContext(ctx).
		Raw("SELECT name, checksum FROM " + migrationsTable).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, r := range rows {
		out[r.Name] = r.Checksum
	}
	return out, nil
}

func listMigrationFiles() ([]string, error) {
	entries, err := fs.ReadDir(FS, "sql")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), ".up.sql") {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)
	return names, nil
}

func checksum(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// splitSQLStatements splits a SQL file into individual statements.
// Handles:
//   - line comments        (-- ... until newline)
//   - block comments       (/* ... */)
//   - single-quoted strings ('...' with ” escape)
//   - double-quoted idents ("...")
//   - dollar-quoted bodies ($$ ... $$, $tag$ ... $tag$)
//
// Empty statements are dropped.
func splitSQLStatements(sql string) []string {
	var stmts []string
	var buf strings.Builder

	var inSingle, inDouble, inLineComment, inBlockComment, inDollar bool
	var dollarTag string

	runes := []rune(sql)
	n := len(runes)

	for i := 0; i < n; i++ {
		c := runes[i]

		// ---- line comment ----
		if inLineComment {
			if c == '\n' {
				inLineComment = false
				buf.WriteRune(c)
			}
			continue
		}

		// ---- block comment ----
		if inBlockComment {
			if c == '*' && i+1 < n && runes[i+1] == '/' {
				inBlockComment = false
				i++ // skip '/'
			}
			continue
		}

		// ---- dollar-quoted ----
		if inDollar {
			// Check if we hit the closing tag.
			if c == '$' {
				remaining := string(runes[i:])
				if strings.HasPrefix(remaining, dollarTag) {
					// Write the full closing tag.
					for j := 0; j < len(dollarTag); j++ {
						buf.WriteRune(runes[i+j])
					}
					i += len(dollarTag) - 1
					inDollar = false
					dollarTag = ""
					continue
				}
			}
			buf.WriteRune(c)
			continue
		}

		// ---- single-quoted string ----
		if inSingle {
			buf.WriteRune(c)
			if c == '\'' {
				// Escaped quote '' stays inside the string.
				if i+1 < n && runes[i+1] == '\'' {
					i++
					buf.WriteRune(runes[i])
				} else {
					inSingle = false
				}
			}
			continue
		}

		// ---- double-quoted identifier ----
		if inDouble {
			buf.WriteRune(c)
			if c == '"' {
				inDouble = false
			}
			continue
		}

		// ---- detect starts ----
		if c == '-' && i+1 < n && runes[i+1] == '-' {
			inLineComment = true
			i++ // skip second '-'
			continue
		}
		if c == '/' && i+1 < n && runes[i+1] == '*' {
			inBlockComment = true
			i++ // skip '*'
			continue
		}
		if c == '\'' {
			inSingle = true
			buf.WriteRune(c)
			continue
		}
		if c == '"' {
			inDouble = true
			buf.WriteRune(c)
			continue
		}
		if c == '$' {
			// Try to read a dollar-quote tag: $tag$ or $$
			j := i + 1
			for j < n && (runes[j] == '_' ||
				(runes[j] >= 'a' && runes[j] <= 'z') ||
				(runes[j] >= 'A' && runes[j] <= 'Z') ||
				(runes[j] >= '0' && runes[j] <= '9')) {
				j++
			}
			if j < n && runes[j] == '$' {
				dollarTag = string(runes[i : j+1])
				buf.WriteString(dollarTag)
				i = j
				inDollar = true
				continue
			}
			// Plain $ — just a character.
			buf.WriteRune(c)
			continue
		}

		// ---- statement delimiter ----
		if c == ';' {
			stmt := strings.TrimSpace(buf.String())
			if stmt != "" {
				stmts = append(stmts, stmt)
			}
			buf.Reset()
			continue
		}

		buf.WriteRune(c)
	}

	// Trailing statement (no semicolon at end of file).
	if stmt := strings.TrimSpace(buf.String()); stmt != "" {
		stmts = append(stmts, stmt)
	}

	return stmts
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
