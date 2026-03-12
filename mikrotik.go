// Package mikrotik provides RouterOS API client
package mikrotik

import (
	"github.com/Butterfly-Student/go-ros/client"
)

// Client is the main Mikrotik client facade
type Client struct {
	conn *client.Client
}

// NewClient creates a new Mikrotik client facade using the provided Config.
func NewClient(cfg client.Config) (*Client, error) {
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: c,
	}, nil
}

// Conn returns the underlying client connection.
func (c *Client) Conn() *client.Client {
	return c.conn
}

// Close closes the underlying connection.
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
