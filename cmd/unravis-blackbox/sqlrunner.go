package main

// sqlrunner.go is the wiring for the separate SQL service process (BBX-5-10).
//
// This code path runs ONLY inside the UnravisBlackboxSql service, as
// NT SERVICE\UnravisBlackboxSql. It writes only to o.Dir (the DACL-locked \sql
// folder), never touches the host recorder's output, opens no connection of
// its own (the runner owns the single connection), serves no HTTP, and reads
// no credentials: the SQL text lives inside internal/bbsql (C1) and nothing
// here passes SQL, parameters or error text onward (C10).
//
// Two things cannot be proven without a real server, so the owner must check
// them on a live instance before this is believed:
//   1. that the SQL service samples a local instance end to end (connect, C8
//      ceiling check, allowlisted query, decode, NDJSON append);
//   2. that host sampling stays unaffected while SQL is slow or backed off.
//      Process separation makes it true by construction here, but the
//      measurable effect needs a loaded server.

import (
	"context"
	"os"
	"time"

	"api-audit-explorer/internal/bbsql"
	"api-audit-explorer/internal/blackbox"
)

// sqlLocalServer is the only SQL target this build can reach. C12 keeps the
// Blackbox to the machine it runs on, and the default instance is the one
// agreed for this release; a named instance needs a probe.server field, which
// config.go deliberately still refuses (owner decision, live check).
const sqlLocalServer = "(local)"

// runSQLServiceLoop is the real deps.sqlRun: it runs the bbsql runner in this
// process until ctx is done. The budgets (C4), the back-off and auto-off (C7)
// and the status file shape (C11) all live inside bbsql.Runner; this function
// only points it at the locked folder and the owner's trust anchors.
func runSQLServiceLoop(ctx context.Context, o sqlRunOptions) error {
	// install-sql sets the folder's ACL; this process only creates it.
	if err := os.MkdirAll(o.Dir, 0o700); err != nil {
		return err
	}
	r := bbsql.NewRunner(o.Dir, sqlConnConfig(o.Config.Probe))
	if o.Log != nil {
		r.Log = o.Log
	}
	return r.Run(ctx)
}

// sqlConnConfig maps the probe block onto the connection settings. Empty
// anchors mean no anchor, and bbsql.Open then fails closed rather than
// accepting an unverifiable certificate. QueryTimeout is the 60 s item's C4
// ceiling; the runner narrows it to 3 s for the 15 s items.
func sqlConnConfig(p blackbox.Probe) bbsql.ConnConfig {
	return bbsql.ConnConfig{
		Server:                sqlLocalServer,
		QueryTimeout:          5 * time.Second,
		ServerCertificate:     p.ServerCertificate,
		Certificate:           p.Certificate,
		HostNameInCertificate: p.HostNameInCertificate,
	}
}
