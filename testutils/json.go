package testutils

import (
	"os"
	"testing"
)

func ReadJSONFromFile(t *testing.T, inputFile string) []byte {
	inputJSON, err := os.ReadFile(inputFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}
	return inputJSON
}
