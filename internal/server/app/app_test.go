package app_test

import (
	"context"
	"testing"

	"github.com/holos-run/holos/internal/server/app"
	"github.com/holos-run/holos/internal/server/testutils"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Run("AppContext", func(t *testing.T) {
		app := app.New()
		assert.NotNil(t, app)

		t.Run("WithContext", func(t *testing.T) {
			app2 := app.WithContext(context.Background())
			assert.NotNil(t, app2)
		})
		t.Run("WithLogger", func(t *testing.T) {
			app2 := app.WithLogger(testutils.TestLogger(t))
			assert.NotNil(t, app2)
		})
		t.Run("ContextLogger", func(t *testing.T) {
			ctx, log := app.ContextLogger()
			assert.NotNil(t, ctx)
			assert.NotNil(t, log)
		})
	})
}
