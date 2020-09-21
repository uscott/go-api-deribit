package api

import (
	"time"
)

// Partially implement Subordinate interface from module go-restarter

type Subordinate struct {
	ChError chan error
	ChDone  chan bool
	Wait    time.Duration
}

func NewSubordinate() *Subordinate {
	s := &Subordinate{}
	s.ChError = make(chan error)
	s.ChDone = make(chan bool)
	return s
}

func (c *Client) Done() chan bool {
	return c.Sub.ChDone
}

func (c *Client) Error() chan error {
	return c.Sub.ChError
}

func (c *Client) WaitTime() time.Duration {
	return c.Sub.Wait
}
