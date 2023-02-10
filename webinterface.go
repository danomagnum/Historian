package main

import (
	"encoding/json"
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
	mux.HandleFunc("/GetWorking/", api_GetWorkingConf)
	mux.HandleFunc("/LoadWorking/", api_LoadWorkingConf)
	mux.HandleFunc("/GetActive/", api_GetActiveConf)
	mux.HandleFunc("/Server/", api_ServerConf)
	mux.HandleFunc("/Providers/", api_ProidersConf)
	mux.Handle("/Providers/CIPClass3/", http.StripPrefix("/Providers/CIPClass3", webApiCIPClass3_Handler))
	addr := fmt.Sprintf("%s:%d", activeConf.General.Host, activeConf.General.Port)
	go http.ListenAndServe(addr, mux)

}

func api_GetActiveConf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	writer := json.NewEncoder(w)
	err := writer.Encode(activeConf)
	if err != nil {
		log.Printf("problem encoding active conf: %v", err)
	}
}

func api_GetWorkingConf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	writer := json.NewEncoder(w)
	err := writer.Encode(workingConf)
	if err != nil {
		log.Printf("problem encoding working conf: %v", err)
	}
}

func api_LoadWorkingConf(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	reader := json.NewDecoder(r.Body)
	err := reader.Decode(&workingConf)
	if err != nil {
		log.Printf("problem decoding working conf: %v", err)
		return
	}
	changes = true
	api_Home(w, r)
}
