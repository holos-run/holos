package holos

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/holos-run/holos/internal/logger"
	"k8s.io/client-go/util/homedir"
)

// NewClientContext loads a ClientContext from the file system if it exists,
// otherwise returns a ClientContext with default values.
func NewClientContext(ctx context.Context) *ClientContext {
	cc := &ClientContext{}
	if cc.Exists() {
		if err := cc.Load(ctx); err != nil {
			logger.FromContext(ctx).WarnContext(ctx, "could not load client context", "err", err)
			return nil
		}
	}
	return cc
}

// ClientContext represents the context the holos api is working in.  Used to
// store and recall values from the filesystem.
type ClientContext struct {
	// OrgID is the organization id of the current context.
	OrgID string `json:"org_id"`
	// UserID is the user id of the current context.
	UserID string `json:"user_id"`
}

func (cc *ClientContext) Save(ctx context.Context) error {
	log := logger.FromContext(ctx)
	config := cc.configFile()
	data, err := json.MarshalIndent(cc, "", "  ")
	if err != nil {
		return err
	}
	if len(data) > 0 {
		data = append(data, '\n')
	}
	if err := os.WriteFile(config, data, 0644); err != nil {
		return err
	}
	log.DebugContext(ctx, "saved", "path", config, "bytes", len(data))
	return nil
}

func (cc *ClientContext) Load(ctx context.Context) error {
	log := logger.FromContext(ctx)
	config := cc.configFile()
	data, err := os.ReadFile(config)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, cc); err != nil {
		return err
	}
	log.DebugContext(ctx, "loaded", "path", config, "bytes", len(data))
	return nil
}

// Exists returns true if the client context file exists.
func (cc *ClientContext) Exists() bool {
	_, err := os.Stat(cc.configFile())
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func (cc *ClientContext) configFile() string {
	config := "client-context.json"
	if home := homedir.HomeDir(); home != "" {
		dir := filepath.Join(home, ".holos")
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			slog.Warn("could not mkdir", "path", dir, "err", err)
		}
		config = filepath.Join(home, ".holos", config)
	}
	return config
}
