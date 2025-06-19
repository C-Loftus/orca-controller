package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewOrcaClient(t *testing.T) {
	client, err := NewOrcaClient()
	require.NoError(t, err)
	defer client.Close()
	require.NotNil(t, client)
	err = client.PresentMessage("test")
	require.NoError(t, err)
}
