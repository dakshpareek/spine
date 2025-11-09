package skeleton

import "testing"

func TestPathForSource(t *testing.T) {
	tests := map[string]string{
		"src/app/service.go":       ".spine/skeletons/src/app/service.skeleton.go",
		"src/app/component.tsx":    ".spine/skeletons/src/app/component.skeleton.tsx",
		"cmd/ctx":                  ".spine/skeletons/cmd/ctx.skeleton",
		"docs/README.md":           ".spine/skeletons/docs/README.skeleton.md",
		"scripts/deploy.sh":        ".spine/skeletons/scripts/deploy.skeleton.sh",
		"src\\windows\\path.go":    ".spine/skeletons/src/windows/path.skeleton.go",
		".github/workflows/ci.yml": ".spine/skeletons/.github/workflows/ci.skeleton.yml",
	}

	for input, expected := range tests {
		if got := PathForSource(input); got != expected {
			t.Fatalf("PathForSource(%q) = %q, expected %q", input, got, expected)
		}
	}
}
