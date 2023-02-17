package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/danomagnum/gologix"
	"github.com/gorilla/mux"
)

type tmplProviderClass3Data struct {
	System System
	Title  string
	Conf   ConfigCIPClass3
}

func cipClass3Init(r *mux.Router) {
	r.HandleFunc("/{name}/Edit/", api_EditCipClass3Conf)
	r.HandleFunc("/{name}/EditEndpoint/", api_EditCipClass3Endpoint)
	r.HandleFunc("/{name}/NewEndpoint/", api_NewCipClass3Endpoint)
	r.HandleFunc("/Add/", api_NewCipClass3Conf)
}

func api_EditCipClass3Conf(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	vars := mux.Vars(r)
	targetName := vars["name"]
	conf, ok := findClass3Endpoint(targetName)
	if !ok {
		log.Printf("Could not find '%s'", targetName)
		return
	}
	editCipClass3Conf(*conf, w, r)
}

func editCipClass3Conf(conf ConfigCIPClass3, w http.ResponseWriter, r *http.Request) {

	dat := tmplProviderClass3Data{
		System: system,
		Title:  fmt.Sprintf("Editing CIP Class 3 Endpoint %s @ %s,%s", conf.Name(), conf.Address, conf.Path),
		Conf:   conf,
	}
	err := templates.ExecuteTemplate(w, "Provider_CIPClass3.html", dat)
	if err != nil {
		log.Printf("problem with template. %v", err)
	}
}

func api_NewCipClass3Conf(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	conf := ConfigCIPClass3{PLCName: "New_CIP_Class3_Endpoint"}
	system.Changes = true

	system.WorkingConfig.DataProviders.CIPClass3 = append(system.WorkingConfig.DataProviders.CIPClass3, conf)

	editCipClass3Conf(conf, w, r)
}

func api_EditCipClass3Endpoint(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	vars := mux.Vars(r)
	targetName := vars["name"]
	conf, ok := findClass3Endpoint(targetName)
	log.Printf("Editing Endpoint %s", r.URL)
	if !ok {
		log.Printf("Could not find '%s'", targetName)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Printf("problem parsing form: %v", err)
		return
	}

	// see what we're editing.
	index_str := r.FormValue("Index")
	index, err := strconv.Atoi(index_str)
	if err != nil {
		log.Printf("invalid endpoint %s not an int: %v", index_str, err)
		return
	} else if index < 0 || index >= len(conf.EndpointList) {
		log.Printf("invalid endpoint %d.  must be 0.. %d", index, len(conf.EndpointList))
		return
	}

	// validate the rate
	rate_str := r.FormValue("Rate")
	rate_str = strings.ReplaceAll(rate_str, " ", "") // get rid of spaces for parsing
	rate, err := time.ParseDuration(rate_str)
	if err != nil {
		log.Printf("invalid rete %s: %v", rate_str, err)
		return
	}

	// validate the type
	ciptype_str := r.FormValue("Type")
	ciptype, err := strconv.Atoi(ciptype_str)
	if err != nil {
		log.Printf("invalid type %s. not an int: %v", ciptype_str, err)
		return
	}

	// load data into that item.
	newendpoint := conf.EndpointList[index]
	newendpoint.Historian = r.FormValue("Historian")
	newendpoint.Name = r.FormValue("Name")
	newendpoint.TagName = r.FormValue("TagName")
	newendpoint.Rate = rate
	newendpoint.TagType = gologix.CIPType(ciptype)
	conf.EndpointList[index] = newendpoint
	system.Changes = true

	editCipClass3Conf(*conf, w, r)

}

func api_NewCipClass3Endpoint(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	vars := mux.Vars(r)
	targetName := vars["name"]
	conf, ok := findClass3Endpoint(targetName)
	log.Printf("Editing Endpoint %s", r.URL)
	if !ok {
		log.Printf("Could not find '%s'", targetName)
		return
	}

	system.Changes = true
	conf.EndpointList = append(conf.EndpointList, EndpointCIPClass3{})

	editCipClass3Conf(*conf, w, r)

}

func findClass3Endpoint(name string) (*ConfigCIPClass3, bool) {
	for i := range system.WorkingConfig.DataProviders.CIPClass3 {
		if system.WorkingConfig.DataProviders.CIPClass3[i].PLCName == name {
			return &system.WorkingConfig.DataProviders.CIPClass3[i], true
		}
	}

	log.Printf("Could not find '%s'", name)
	return &ConfigCIPClass3{}, false
}
