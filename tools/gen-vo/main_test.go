package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateVO(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "genvo-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
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
	if err := os.WriteFile(srcFile, []byte(dummySrc), 0644); err != nil {
		t.Fatalf("failed to write dummy source: %v", err)
	}

	// Execute generate function (to be implemented)
	err = runGenerator(srcFile)
	if err != nil {
		t.Fatalf("runGenerator returned error: %v", err)
	}

	genFile := filepath.Join(tmpDir, "dummy_gen.go")
	content, err := os.ReadFile(genFile)
	if err != nil {
		t.Fatalf("failed to read generated file: %v", err)
	}

	expectedMethods := []string{
		"func (id DummyID) Value() string",
		"func (id DummyID) Equals(other DummyID) bool",
		"func (id DummyID) String() string",
		"func (id GroupedID) Value() string",
		"func (id GroupedID) Equals(other GroupedID) bool",
		"func (id GroupedID) String() string",
	}

	for _, method := range expectedMethods {
		if !strings.Contains(string(content), method) {
			t.Errorf("missing generated method: %q", method)
		}
	}

	unexpectedMethods := []string{
		"func (id InvalidID)",
		"InvalidID",
	}

	for _, method := range unexpectedMethods {
		if strings.Contains(string(content), method) {
			t.Errorf("unexpected method generated for non-string VO: %q", method)
		}
	}
}
