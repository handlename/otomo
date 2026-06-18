package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateVO(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "genvo-test")
	require.NoError(t, err, "failed to create temp dir")
	defer os.RemoveAll(tmpDir)

	dummySrc := `package testpkg
// @vo
type DummyID struct {
	value string
}

type (
	// @vo
	GroupedID struct {
		value string
	}

	// @vo
	InvalidID struct {
		value int
	}
)
`
	srcFile := filepath.Join(tmpDir, "dummy.go")
	err = os.WriteFile(srcFile, []byte(dummySrc), 0644)
	require.NoError(t, err, "failed to write dummy source")

	err = runGenerator(srcFile)
	require.NoError(t, err, "runGenerator returned error")

	genFile := filepath.Join(tmpDir, "dummy_gen.go")
	contentBytes, err := os.ReadFile(genFile)
	require.NoError(t, err, "failed to read generated file")
	content := string(contentBytes)

	expectedMethods := []string{
		"func (id DummyID) Value() string",
		"func (id DummyID) Equals(other DummyID) bool",
		"func (id DummyID) String() string",
		"func (id GroupedID) Value() string",
		"func (id GroupedID) Equals(other GroupedID) bool",
		"func (id GroupedID) String() string",
	}

	for _, method := range expectedMethods {
		assert.Contains(t, content, method, "missing generated method")
	}

	unexpectedMethods := []string{
		"func (id InvalidID)",
		"InvalidID",
	}

	for _, method := range unexpectedMethods {
		assert.NotContains(t, content, method, "unexpected method generated for non-string VO")
	}
}
