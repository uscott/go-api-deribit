package test

import (
	"testing"

	"github.com/uscott/go-api-deribit/api"
)

func TestClientCreate(t *testing.T) {
	_, err := api.New(api.DfltCnfg())
	if err != nil {
		t.Fatal(err.Error())
	}
}
