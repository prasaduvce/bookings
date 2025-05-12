package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/prasaduvce/bookings/internal/config"
	"github.com/prasaduvce/bookings/internal/models"
)

var appConfig *config.AppConfig
var pathToTemplates = "./templates"

func NewRenderer(a *config.AppConfig) {
	appConfig = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = appConfig.Session.PopString(r.Context(), "flash")
	td.Error = appConfig.Session.PopString(r.Context(), "error")
	td.Warning = appConfig.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderHtml is a function that renders HTML templates.
func Template(w http.ResponseWriter, templatePath string, td *models.TemplateData, r *http.Request) error {

	//var tc config.AppConfig

	//get AppConfig template cache
	

	//create a template cache
	//tc, err := a.CreateTemplateCache()

	var tc map[string]*template.Template
	if appConfig.UseCache {
		tc = appConfig.Templates
		//log.Println("Using cached template")
	} else {
		tc, _ = CreateTemplateCache()
		//log.Println("Creating the cahced template")
	}

	t, ok := tc[templatePath]
	if !ok {
		log.Println("Could not get template from cache")
		return errors.New("could not get template from cache")
	}
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println("Error executing template: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to response: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	//get requested template from cache

	//render the template

	/*
	prasedTemplate, err := template.ParseFiles("./templates/" + templatePath, "./templates/base.layout.tmpl")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	errTemplate := prasedTemplate.Execute(w, nil)
	if errTemplate != nil {
		fmt.Println("Error executing template: ", errTemplate)
		return
	}*/
	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Templat
	myCache := map[string]*template.Template{}

	//get all the files names *.page.tmpl from templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	//range through all the files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		//parse the page and base layout
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Println("Error parsing template: ", err)
			return myCache, err
		}

		//parse the base layout
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			log.Println("Error parsing base layout: ", err)
			return myCache, err
		}
		//parse the partials
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				log.Println("Error parsing base layout: ", err)
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}

