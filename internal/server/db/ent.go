// Package db manages database client and schema migrations to interact with entities.
package db

import (
	"context"
	"database/sql"
	"time"

	"entgo.io/ent/dialect"
	"github.com/holos-run/holos/internal/ent"
	"github.com/holos-run/holos/pkg/errors"
	"github.com/holos-run/holos/pkg/holos"
)

// Conn holds database connection info
type Conn struct {
	Client *ent.Client
	DB     *sql.DB
	Driver dialect.Driver
}

func Client(cfg *holos.Config) (Conn, error) {
	var clientFactory ClientFactory
	if cfg.ServerConfig == nil || cfg.ServerConfig.DatabaseURI() == "" {
		clientFactory = NewMemoryClientFactory(cfg)
	} else {
		clientFactory = NewPGXClientFactory(cfg)
	}
	return clientFactory.New()
}

type ClientFactory interface {
	New() (Conn, error)
}

// withHooks adds our custom hooks to the database client.
func withHooks(client *ent.Client) *ent.Client {
	client.Use(setUpdatedToCreatedOnCreate)
	return client
}

// setUpdatedToCreatedOnCreate ensures the updatedAt field matches the createdAt field on creation.
func setUpdatedToCreatedOnCreate(next ent.Mutator) ent.Mutator {
	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		if createdAt, ok := m.Field("created_at"); ok {
			if _, ok := m.Field("updated_at"); ok && m.Op().Is(ent.OpCreate) {
				err := m.SetField("updated_at", createdAt.(time.Time))
				if err != nil {
					return nil, errors.Wrap(err)
				}
			}
		}
		return next.Mutate(ctx, m)
	})
}
