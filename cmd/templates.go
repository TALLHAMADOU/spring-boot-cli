package cmd

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// renderTemplate parses and executes a named Java source template with the given data.
func renderTemplate(name string, data any) (string, error) {
	tmplPath := fmt.Sprintf("templates/%s.tmpl", name)
	tmplContent, err := templateFS.ReadFile(tmplPath)
	if err != nil {
		return "", fmt.Errorf("template %q not found: %w", name, err)
	}

	tmpl, err := template.New(name).Parse(string(tmplContent))
	if err != nil {
		return "", fmt.Errorf("parse template %q: %w", name, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template %q: %w", name, err)
	}
	return buf.String(), nil
}
