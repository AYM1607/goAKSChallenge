package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	invalidDir       = "testdata/invalid"
	invalidSchemaDir = "testdata/invalidSchema"
	validDir         = "testdata/valid"
)

// TODO: Convert this set of tests to a table test to delete dupe code.

func TestInvalid(t *testing.T) {
	fis, err := ioutil.ReadDir(invalidDir)
	require.NoError(t, err, "Test dir for invalid records doesn't exist")
	for _, fi := range fis {
		data, err := os.ReadFile(filepath.Join(invalidDir, fi.Name()))
		require.NoError(t, err)
		_, err = newRecord(data)
		require.Equal(t, ErrUnparsable, err, "Invalid files shouldn't be able to be parsed")
	}
}

func TestInvalidSchema(t *testing.T) {
	fis, err := ioutil.ReadDir(invalidSchemaDir)
	require.NoError(t, err, "Test dir for invalid schema records doesn't exist")
	for _, fi := range fis {
		data, err := os.ReadFile(filepath.Join(invalidSchemaDir, fi.Name()))
		require.NoError(t, err)
		_, err = newRecord(data)
		require.Error(t, err, "Record creation should field if any field is invalid")
		require.NotEqual(t, ErrUnparsable, err, "If invalid fields were found, the err should refelct that.")
	}
}

func TestValid(t *testing.T) {
	fis, err := ioutil.ReadDir(validDir)
	require.NoError(t, err, "Test dir for valid records doesn't exist")
	for _, fi := range fis {
		data, err := os.ReadFile(filepath.Join(validDir, fi.Name()))
		require.NoError(t, err)
		_, err = newRecord(data)
		require.NoError(t, err, "Valid files should be parsed correctly")
	}
}
