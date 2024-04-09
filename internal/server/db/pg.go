package db

import (
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/holos-run/holos/internal/server/app"
	"github.com/holos-run/holos/internal/server/ent"
	"github.com/holos-run/holos/pkg/errors"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPGXClientFactory returns a PGXClientFactory implementation of ClientFactory
func NewPGXClientFactory(app app.App) *PGXClientFactory {
	return &PGXClientFactory{app: app}
}

// PGXClientFactory produces pgx clients suitable for live workloads
type PGXClientFactory struct {
	app app.App
}

// New returns a new ent.Client using pgx with PostgreSQL
func (mc *PGXClientFactory) New() (Conn, error) {
	uri := mc.app.Config.DatabaseURI()
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return Conn{}, errors.Wrap(fmt.Errorf("could not open pgx: %w", err))
	}
	drv := entsql.OpenDB(dialect.Postgres, db)
	client := withHooks(ent.NewClient(ent.Driver(drv)))
	return Conn{client, db, drv}, nil
}
