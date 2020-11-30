package test

import (
	"testing"

	"github.com/uscott/go-api-deribit/api"
)

func newClient() (*api.Client, error) {
	return api.New(api.DfltCnfg())
}

func TestClientCreate(t *testing.T) {
	var err error
	_, err = newClient()
	if err != nil {
		t.Fatal(err.Error())
	}
}
