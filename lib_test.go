// Copyright 2025 Colton Loftus
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpeechFunctions(t *testing.T) {
	client, err := NewOrcaClient()
	require.NoError(t, err)
	defer client.Close()

	t.Run("TestChangePitch", func(t *testing.T) {
		err = client.SpeechAndVerbosityManager.InterruptSpeech(true)
		require.NoError(t, err)

		for i := 0; i < 5; i++ {
			err = client.PresentMessage("test")
			require.NoError(t, err)
			err = client.SpeechAndVerbosityManager.DecreasePitch(true)
			require.NoError(t, err)
		}
		for i := 0; i < 5; i++ {
			err = client.PresentMessage("test")
			require.NoError(t, err)
			err = client.SpeechAndVerbosityManager.IncreasePitch(true)
			require.NoError(t, err)
		}
	})

	t.Run("TestSpeak", func(t *testing.T) {
		err = client.SpeechAndVerbosityManager.InterruptSpeech(true)
		require.NoError(t, err)
		err = client.PresentMessage("This is a message from go")
		require.NoError(t, err)
	})
}
