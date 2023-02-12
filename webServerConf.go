package main

import (
	"html/template"
	"log"
	"net/http"
)

type tmplServerConfData struct {
	System System
	Title  string
	Conf   ConfigGeneral
}

func api_ServerConf(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	dat := tmplServerConfData{
		System: system,
		Title:  "ServerConf",
		Conf:   system.WorkingConfig.General,
	}
	err := templates.ExecuteTemplate(w, "Server.html", dat)
	if err != nil {
		log.Printf("problem with template. %v", err)
	}
}
