package main

import (
	"html/template"
	"log"
	"net/http"
)

type tmplProvidersData struct {
	Title string
	Conf  ConfigDataProviders
}

func api_ProidersConf(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob("./templates/*") // TODO: remove once page debug is done
	dat := tmplProvidersData{
		Title: "ServerConf",
		Conf:  workingConf.DataProviders,
	}
	err := templates.ExecuteTemplate(w, "Providers.html", dat)
	if err != nil {
		log.Printf("problem with template. %v", err)
	}
}
