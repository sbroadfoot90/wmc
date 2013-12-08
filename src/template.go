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
		"index",
		"restaurant",
		"editRestaurant",
		"top10",
	}
	fm := template.FuncMap{
		"UserName":       userName,
		"RestaurantName": restaurantName,
		"eq":             equals,
		"neq":            notequals,
		"seq":             sequals,
	}
	for _, templateName := range templateNames {
		root := filepath.Join(templatePath, "root.tmpl")
		templates[templateName] = template.Must(template.New("").Funcs(fm).ParseFiles(root, filepath.Join(templatePath, templateName+".tmpl")))
	}
}

func equals(a, b int) bool {
	return a == b
}

func sequals(a, b string) bool {
	return a == b
}

func notequals(a, b int) bool {
	return a != b
}
