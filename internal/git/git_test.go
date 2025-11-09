package git

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

type expectedCommand struct {
	name   string
	args   []string
	output []byte
	err    error
}

type fakeRunner struct {
	t     *testing.T
	cmds  []expectedCommand
	index int
}

func newFakeRunner(t *testing.T, cmds []expectedCommand) *fakeRunner {
	return &fakeRunner{t: t, cmds: cmds}
}

func (f *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	f.t.Helper()

	if f.index >= len(f.cmds) {
		f.t.Fatalf("unexpected command: %s %v", name, args)
	}

	expected := f.cmds[f.index]
	f.index++

	if expected.name != name {
		f.t.Fatalf("expected command %s, got %s", expected.name, name)
	}
	if !reflect.DeepEqual(expected.args, args) {
		f.t.Fatalf("expected args %v, got %v", expected.args, args)
	}

	return expected.output, expected.err
}

func (f *fakeRunner) assertAllCommandsUsed() {
	if f.index != len(f.cmds) {
		f.t.Fatalf("not all expected commands were used: %d remaining", len(f.cmds)-f.index)
	}
}

func resetRunner(t *testing.T) {
	t.Helper()
	runner = execCommandRunner{}
}

func TestIsGitRepoTrue(t *testing.T) {
	fake := newFakeRunner(t, []expectedCommand{
		{
			name:   "git",
			args:   []string{"rev-parse", "--is-inside-work-tree"},
			output: []byte("true\n"),
		},
	})
	runner = fake
	t.Cleanup(func() {
		fake.assertAllCommandsUsed()
		resetRunner(t)
	})

	if !IsGitRepo() {
		t.Fatalf("expected repository detection to succeed")
	}
}

func TestIsGitRepoFalse(t *testing.T) {
	fake := newFakeRunner(t, []expectedCommand{
		{
			name: "git",
			args: []string{"rev-parse", "--is-inside-work-tree"},
			err:  errors.New("fatal: not a git repository"),
		},
	})
	runner = fake
	t.Cleanup(func() {
		fake.assertAllCommandsUsed()
		resetRunner(t)
	})

	if IsGitRepo() {
		t.Fatalf("expected repository detection to return false")
	}
}

func TestGetModifiedFiles(t *testing.T) {
	fake := newFakeRunner(t, []expectedCommand{
		{
			name:   "git",
			args:   []string{"rev-parse", "--is-inside-work-tree"},
			output: []byte("true\n"),
		},
		{
			name:   "git",
			args:   []string{"diff", "--name-only", "HEAD"},
			output: []byte("file1.go\nnested/file2.ts\n"),
		},
	})
	runner = fake
	t.Cleanup(func() {
		fake.assertAllCommandsUsed()
		resetRunner(t)
	})

	files, err := GetModifiedFiles()
	if err != nil {
		t.Fatalf("GetModifiedFiles error: %v", err)
	}

	expected := []string{"file1.go", "nested/file2.ts"}
	if !reflect.DeepEqual(files, expected) {
		t.Fatalf("expected files %v, got %v", expected, files)
	}
}

func TestGetModifiedFilesNotGit(t *testing.T) {
	fake := newFakeRunner(t, []expectedCommand{
		{
			name: "git",
			args: []string{"rev-parse", "--is-inside-work-tree"},
			err:  errors.New("fatal: not a git repository"),
		},
	})
	runner = fake
	t.Cleanup(func() {
		fake.assertAllCommandsUsed()
		resetRunner(t)
	})

	if _, err := GetModifiedFiles(); !errors.Is(err, ErrNotGit) {
		t.Fatalf("expected ErrNotGit, got %v", err)
	}
}

func TestGetUntrackedFiles(t *testing.T) {
	fake := newFakeRunner(t, []expectedCommand{
		{
			name:   "git",
			args:   []string{"rev-parse", "--is-inside-work-tree"},
			output: []byte("true\n"),
		},
		{
			name:   "git",
			args:   []string{"ls-files", "--others", "--exclude-standard"},
			output: []byte("new.ts\nnested/new.go\n"),
		},
	})
	runner = fake
	t.Cleanup(func() {
		fake.assertAllCommandsUsed()
		resetRunner(t)
	})

	files, err := GetUntrackedFiles()
	if err != nil {
		t.Fatalf("GetUntrackedFiles error: %v", err)
	}

	expected := []string{"new.ts", "nested/new.go"}
	if !reflect.DeepEqual(files, expected) {
		t.Fatalf("expected files %v, got %v", expected, files)
	}
}

func TestGetUntrackedFilesNotGit(t *testing.T) {
	fake := newFakeRunner(t, []expectedCommand{
		{
			name: "git",
			args: []string{"rev-parse", "--is-inside-work-tree"},
			err:  errors.New("fatal: not a git repository"),
		},
	})
	runner = fake
	t.Cleanup(func() {
		fake.assertAllCommandsUsed()
		resetRunner(t)
	})

	if _, err := GetUntrackedFiles(); !errors.Is(err, ErrNotGit) {
		t.Fatalf("expected ErrNotGit, got %v", err)
	}
}

func TestGetModifiedFilesFallback(t *testing.T) {
	root := t.TempDir()

	oldPath := filepath.Join(root, "old.txt")
	if err := os.WriteFile(oldPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("failed to write old file: %v", err)
	}
	oldTime := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(oldPath, oldTime, oldTime); err != nil {
		t.Fatalf("failed to set old file time: %v", err)
	}

	newPath := filepath.Join(root, "new.txt")
	if err := os.WriteFile(newPath, []byte("new"), 0o644); err != nil {
		t.Fatalf("failed to write new file: %v", err)
	}

	since := time.Now().Add(-30 * time.Minute)
	files, err := GetModifiedFilesFallback(root, since)
	if err != nil {
		t.Fatalf("GetModifiedFilesFallback error: %v", err)
	}

	expected := []string{"new.txt"}
	if !reflect.DeepEqual(files, expected) {
		t.Fatalf("expected files %v, got %v", expected, files)
	}
}
