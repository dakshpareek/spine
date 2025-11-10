package skeleton

import "testing"

func TestPathForSource(t *testing.T) {
	tests := map[string]string{
		"src/app/service.go":       ".ctx/skeletons/src/app/service.skeleton.go",
		"src/app/component.tsx":    ".ctx/skeletons/src/app/component.skeleton.tsx",
		"cmd/ctx":                  ".ctx/skeletons/cmd/ctx.skeleton",
		"docs/README.md":           ".ctx/skeletons/docs/README.skeleton.md",
		"scripts/deploy.sh":        ".ctx/skeletons/scripts/deploy.skeleton.sh",
		"src\\windows\\path.go":    ".ctx/skeletons/src/windows/path.skeleton.go",
		".github/workflows/ci.yml": ".ctx/skeletons/.github/workflows/ci.skeleton.yml",
	}

	for input, expected := range tests {
		if got := PathForSource(input); got != expected {
			t.Fatalf("PathForSource(%q) = %q, expected %q", input, got, expected)
		}
	}
}
