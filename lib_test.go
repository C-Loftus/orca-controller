// Copyright 2025 Colton Loftus
// SPDX-License-Identifier: AGPL-3.0-only

package pkg

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

		for range 5 {
			err = client.PresentMessage("test")
			require.NoError(t, err)
			err = client.SpeechAndVerbosityManager.DecreasePitch(true)
			require.NoError(t, err)
		}
		for range 5 {
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

func TestVersion(t *testing.T) {
	client, err := NewOrcaClient()
	require.NoError(t, err)
	defer client.Close()
	version, err := client.GetVersion()
	require.NoError(t, err)
	require.NotEmpty(t, version)
}

func TestGetters(t *testing.T) {
	client, err := NewOrcaClient()
	require.NoError(t, err)
	defer client.Close()
	rate, err := client.SpeechAndVerbosityManager.Rate()
	require.NoError(t, err)
	require.Greater(t, rate, float64(0))
}

func TestGetThenSet(t *testing.T) {
	client, err := NewOrcaClient()
	require.NoError(t, err)
	defer client.Close()

	rate, err := client.SpeechAndVerbosityManager.Rate()
	require.NoError(t, err)
	require.NotEmpty(t, rate)

	err = client.SpeechAndVerbosityManager.SetRate(rate)
	require.NoError(t, err)
}
