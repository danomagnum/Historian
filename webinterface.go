package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"html/template"

	"github.com/gorilla/mux"
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
	//router := http.NewServeMux()
	router := mux.NewRouter()
	fs := http.FileServer(http.Dir("./static"))
	router.HandleFunc("/", api_Home)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	router.HandleFunc("/GetWorking/", api_GetWorkingConf)
	router.HandleFunc("/LoadWorking/", api_LoadWorkingConf)
	router.HandleFunc("/ApplyWorking/", api_ApplyWorkingConf)
	router.HandleFunc("/GetActive/", api_GetActiveConf)
	router.HandleFunc("/Server/", api_ServerConf)
	router.HandleFunc("/Providers/", api_ProidersConf)
	//router.PathPrefix("/Providers/CIPClass3/").Handler(http.StripPrefix("/Providers/CIPClass3", webApiCIPClass3_Handler))
	cipClass3Init(router.PathPrefix("/Providers/CIPClass3").Subrouter())
	addr := fmt.Sprintf("%s:%d", activeConf.General.Host, activeConf.General.Port)
	go http.ListenAndServe(addr, router)
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

func api_ApplyWorkingConf(w http.ResponseWriter, r *http.Request) {
	activeContextCancel()
}
