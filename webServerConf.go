package main

import (
	"html/template"
	"log"
	"net/http"
)

type tmplServerConfData struct {
	Changes bool
	Title   string
	Conf    ConfigGeneral
}

func api_ServerConf(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob("./templates/*") // TODO: remove once page debug is done
	dat := tmplServerConfData{
		Changes: changes,
		Title:   "ServerConf",
		Conf:    workingConf.General,
	}
	err := templates.ExecuteTemplate(w, "Server.html", dat)
	if err != nil {
		log.Printf("problem with template. %v", err)
	}
}
