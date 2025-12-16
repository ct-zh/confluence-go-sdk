package template

import "embed"

// content holds our static web server content.
//
//go:embed *.html
var content embed.FS

// Loader defines the interface for loading template content
type Loader interface {
	Load(name string) (string, error)
}

// BuiltinLoader loads templates from the embedded filesystem
type BuiltinLoader struct{}

// NewBuiltinLoader creates a new BuiltinLoader
func NewBuiltinLoader() *BuiltinLoader {
	return &BuiltinLoader{}
}

// Load retrieves the content of a template from the embedded FS
// name should be the filename without path, e.g., "swagger.html"
func (l *BuiltinLoader) Load(name string) (string, error) {
	data, err := content.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
