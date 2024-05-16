// Package client provides client configuration for the holos cli.
package client

import (
	"flag"

	"github.com/holos-run/holos/internal/holos"
	"github.com/holos-run/holos/internal/token"
)

type Config struct {
	holos  *holos.Config
	client *holos.ClientConfig
	token  *token.Config
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

func NewConfig(cfg *holos.Config) *Config {
	return &Config{
		holos:  cfg,
		client: holos.NewClientConfig(),
		token:  token.NewConfig(),
	}
}
