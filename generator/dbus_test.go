package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetModules(t *testing.T) {
	modules, err := get_modules()
	require.NoError(t, err)
	require.NotEmpty(t, modules)
	require.Contains(t, modules, "SpeechAndVerbosityManager")
}

func TestGetCommands(t *testing.T) {
	module := "SpeechAndVerbosityManager"
	commands, err := get_commands_for_module(module)
	require.NoError(t, err)
	require.NotEmpty(t, commands)
	for _, cmd := range commands {
		const tupleLength = 2
		require.Len(t, cmd, tupleLength)

	}
}
