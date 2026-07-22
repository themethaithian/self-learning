// Package migrations embeds the SQL migration files so the built API binary
// is self-contained and needs no filesystem access to migrate the database.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
