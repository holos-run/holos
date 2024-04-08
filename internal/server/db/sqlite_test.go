package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryClientFactory(t *testing.T) {
	t.Run("MemoryClientFactory", func(t *testing.T) {
		mcf := MemoryClientFactory{}
		conn, err := mcf.New()
		assert.NoError(t, err)
		client := conn.Client
		assert.NoError(t, client.Schema.Create(context.Background()))

		// Create something
		t.Run("CreateUser", func(t *testing.T) {
			uc := client.User.Create().
				SetName("Foo").
				SetEmail("foo@example.com")
			_, err := uc.Save(context.Background())
			assert.NoError(t, err)
		})
	})
}
