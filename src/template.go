package wmc

import (
	"appengine"
	"time"
	"strconv"
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
		"FormatDate":formatDate,
		"eq":             equals,
		"neq":            notequals,
	}
	for _, templateName := range templateNames {
		root := filepath.Join(templatePath, "root.tmpl")
		templates[templateName] = template.Must(template.New("").Funcs(fm).ParseFiles(root, filepath.Join(templatePath, templateName+".tmpl")))
	}
}

func userName(c appengine.Context, id string) string {
	p := retrieveProfile(c, id)
	if p == nil {
		return ""
	}
	return p.Name
}

func restaurantName(c appengine.Context, rid string) string {
	rest := retrieveRestaurant(c, rid)
	if rest == nil {
		return ""
	}
	return rest.Name
}

func formatDate(t time.Time) string {
	d := time.Since(t)
	var v int
	var typ string
	if d.Minutes() < 1 {
		v = int(d.Seconds())
		typ = "second"
	} else if d.Hours() < 1 {
		v = int(d.Minutes())
		typ = "minute"
	} else if d.Hours() < 24 {
		v = int(d.Hours())
		typ = "hour"
	} else {
		v = int(d.Hours()/24)
		typ = "day"
	}
	
	if v != 1 {
		typ = typ + "s"
	}
	
	return strconv.Itoa(v) + " " + typ + " ago"
}

func equals(a, b int) bool {
	return a == b
}

func notequals(a, b int) bool {
	return a != b
}
