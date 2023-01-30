package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"html/template"
)

const templatedir = "templates/"

//func init() {
//mime.AddExtensionType(".css", "text/css")
//}

var templates *template.Template

func WebAPIStart() {
	var err error
	templates, err = template.ParseGlob("./templates/*")
	if err != nil {
		log.Println("Cannot parse templates:", err)
		os.Exit(-1)
	}
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.HandleFunc("/", api_Home)
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/Server/", api_ServerConf)
	mux.HandleFunc("/Providers/", api_ProidersConf)
	addr := fmt.Sprintf("%s:%d", activeConf.General.Host, activeConf.General.Port)
	go http.ListenAndServe(addr, mux)

}
