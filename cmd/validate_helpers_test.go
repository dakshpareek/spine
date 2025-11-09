package cmd

import "testing"

func TestValidateStrictFailure(t *testing.T) {
	dir := t.TempDir()
	writeTempFile(t, dir, "main.go", "package main\n")

	_, _ = executeCommand(t, dir, "init")

	cleanup := changeDir(t, dir)
	defer cleanup()

	err := runValidate(validateOptions{strict: true})
	if err == nil {
		t.Fatalf("expected error in strict mode when issues exist")
	}
}
