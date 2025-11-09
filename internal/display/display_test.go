package display

import "testing"

func TestFormats(t *testing.T) {
	if Success("ok") == "" {
		t.Fatalf("expected success output")
	}
	if Warning("warn") == "" {
		t.Fatalf("expected warning output")
	}
	if Info("info") == "" {
		t.Fatalf("expected info output")
	}
	if Error("error") == "" {
		t.Fatalf("expected error output")
	}
	if Bold("bold") == "" {
		t.Fatalf("expected bold output")
	}
}
