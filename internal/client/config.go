// Package client provides client configuration for the holos cli.
package client

import (
	"context"
	"flag"

	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/token"
)

func NewConfig(cfg *holos.Config) *Config {
	return &Config{
		holos:   cfg,
		client:  holos.NewClientConfig(),
		context: holos.NewClientContext(context.Background()),
		token:   token.NewConfig(),
	}
}

type Config struct {
	holos   *holos.Config
	client  *holos.ClientConfig
	context *holos.ClientContext
	token   *token.Config
}

func (c *Config) ClientFlagSet() *flag.FlagSet {
	if c == nil {
		return nil
	}
	return c.client.FlagSet()
}

func (c *Config) TokenFlagSet() *flag.FlagSet {
	if c == nil {
		return nil
	}
	return c.token.FlagSet()
}

func (c *Config) Token() *token.Config {
	if c == nil {
		return nil
	}
	return c.token
}

func (c *Config) Client() *holos.ClientConfig {
	if c == nil {
		return nil
	}
	return c.client
}

// Context returns the ClientContext useful to get the OrgID and UserID for rpc
// calls.
func (c *Config) Context() *holos.ClientContext {
	if c == nil {
		return nil
	}
	return c.context
}

// Holos returns the *holos.Config
func (c *Config) Holos() *holos.Config {
	if c == nil {
		return nil
	}
	return c.holos
}
