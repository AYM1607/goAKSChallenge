package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var invalidFiles = []string{
	"./testData/invalid.yaml",
}

var invalidSchemaFiles = []string{
	"./testData/missingVersion.yaml",
	"./testData/invalidEmail.yaml",
}

var validFiles = []string{
	"./testData/valid1.yaml",
	"./testData/valid2.yaml",
}

func TestInvalid(t *testing.T) {
	for _, f := range invalidFiles {
		data, err := os.ReadFile(f)
		require.NoError(t, err)
		_, err = newRecord(data)
		require.Equal(t, ErrUnparsable, err, "Invalid files shouldn't be able to be parsed")
	}
}

func TestInvalidSchema(t *testing.T) {
	for _, f := range invalidSchemaFiles {
		data, err := os.ReadFile(f)
		require.NoError(t, err)
		_, err = newRecord(data)
		require.Error(t, err, "Record creation should field if any field is invalid")
		require.NotEqual(t, ErrUnparsable, err, "If invalid fields were found, the err should refelct that.")
	}
}

func TestValid(t *testing.T) {
	for _, f := range validFiles {
		data, err := os.ReadFile(f)
		require.NoError(t, err)
		_, err = newRecord(data)
		require.NoError(t, err, "Valid files should be parsed correctly")
	}
}
