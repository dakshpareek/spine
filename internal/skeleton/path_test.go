package skeleton

import "testing"

func TestPathForSource(t *testing.T) {
	tests := map[string]string{
		"src/app/service.go":       ".code-context/skeletons/src/app/service.skeleton.go",
		"src/app/component.tsx":    ".code-context/skeletons/src/app/component.skeleton.tsx",
		"cmd/ctx":                  ".code-context/skeletons/cmd/ctx.skeleton",
		"docs/README.md":           ".code-context/skeletons/docs/README.skeleton.md",
		"scripts/deploy.sh":        ".code-context/skeletons/scripts/deploy.skeleton.sh",
		"src\\windows\\path.go":    ".code-context/skeletons/src/windows/path.skeleton.go",
		".github/workflows/ci.yml": ".code-context/skeletons/.github/workflows/ci.skeleton.yml",
	}

	for input, expected := range tests {
		if got := PathForSource(input); got != expected {
			t.Fatalf("PathForSource(%q) = %q, expected %q", input, got, expected)
		}
	}
}
