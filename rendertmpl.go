package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var tc = make(map[string]*template.Template)

func RenderTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// check if we already have the template
	_, inMap := tc[t]
	if !inMap { // we don't have it in cache
		// we need to create a new template
		fmt.Println("Creating a new template and adding it to cache for further usage")
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		// we have it in the cache
		fmt.Println("Using cached templates")
	}
	tmpl = tc[t]
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

// create a cahch for templates
func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/index.html",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}
	// add template to cache
	tc[t] = tmpl
	return nil
}
