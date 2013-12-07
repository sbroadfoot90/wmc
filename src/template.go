package wmc

import (
	"html/template"
	"path/filepath"
)

var templates map[string]*template.Template = make(map[string]*template.Template)

func importTemplates(templatePath string) {
	templateNames := []string{
		"profile",
		"edit",
		"indexLoggedIn",
		"index",
	}
	
	for _, templateName := range templateNames {
		root := filepath.Join(templatePath, "root.tmpl")
		templates[templateName] = template.Must(template.New("").ParseFiles(root, filepath.Join(templatePath, templateName+".tmpl")))
	}
}