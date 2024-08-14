package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestGenerateNonce checks if the nonce generated is of the correct length and is composed of allowed characters.
func TestGenerateNonce(t *testing.T) {
	length := 16
	nonce := GenerateNonce(length)

	require.Equal(t, len(nonce), length)

	for _, char := range nonce {
		require.True(t, isCharInAlphabet(char))
	}
}

// Helper function to check if a character is part of the allowed alphabet
func isCharInAlphabet(char rune) bool {
	for _, a := range alphabet {
		if a == char {
			return true
		}
	}
	return false
}

// TestGenerateSignature checks if the generated signature is of the correct length and consistency.
func TestGenerateSignature(t *testing.T) {
	str := "test string"
	expectedLength := 64 // sha256 hash in hex is 64 characters long
	sig := GenerateSignature(str)

	require.Equal(t, len(sig), expectedLength)

	// Test with a known value
	knownStr := "hello"
	expectedSig := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	actualSig := GenerateSignature(knownStr)

	require.Equal(t, actualSig, expectedSig)
}
