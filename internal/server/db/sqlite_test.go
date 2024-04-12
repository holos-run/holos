package db

import (
	"context"
	"testing"

	"github.com/holos-run/holos/internal/server/testutils"
	"github.com/holos-run/holos/pkg/holos"
	"github.com/stretchr/testify/assert"
)

func TestMemoryClientFactory(t *testing.T) {
	t.Run("MemoryClientFactory", func(t *testing.T) {
		cfg := holos.New(holos.Logger(testutils.TestLogger(t)))
		mcf := MemoryClientFactory{cfg: cfg}
		conn, err := mcf.New()
		assert.NoError(t, err)
		client := conn.Client
		assert.NoError(t, client.Schema.Create(context.Background()))

		// Create something
		t.Run("CreateUser", func(t *testing.T) {
			uc := client.User.Create().
				SetName("Foo").
				SetIss("https://login.example.com").
				SetSub("1234567890").
				SetEmail("foo@example.com")
			_, err := uc.Save(context.Background())
			assert.NoError(t, err)
		})
	})
}
