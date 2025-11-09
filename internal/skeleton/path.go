package skeleton

import (
	"path/filepath"
	"strings"
)

const (
	// DirRoot is the root directory used to store skeleton files.
	DirRoot = ".spine/skeletons"
)

// PathForSource returns the relative skeleton path corresponding to a source file path.
func PathForSource(source string) string {
	source = strings.ReplaceAll(source, "\\", "/")
	source = filepath.ToSlash(source)
	ext := filepath.Ext(source)
	base := strings.TrimSuffix(source, ext)

	var builder strings.Builder
	builder.WriteString(DirRoot)
	builder.WriteRune('/')
	builder.WriteString(base)
	builder.WriteString(".skeleton")
	builder.WriteString(ext)

	return builder.String()
}
