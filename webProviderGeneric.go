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

type DataProvider interface {
	Name() string
	Config() any
	SetConfig(any)
}

type apiGenericConfig[T any] struct {
	ConfTypeName string
	Confs        []T
}

type tmplProviderGenericData[T any] struct {
	System System
	Title  string
	Conf   T
}

func (gc apiGenericConfig[T]) Init(r *mux.Router) {
	r.HandleFunc("/{name}/Edit/", gc.api_EditConf)
	r.HandleFunc("/{name}/EditEndpoint/", gc.api_EditEndpoint)
	r.HandleFunc("/{name}/NewEndpoint/", gc.api_NewEndpoint)
	r.HandleFunc("/Add/", gc.api_NewConf)
}

func (gc apiGenericConfig[T]) api_EditConf(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	vars := mux.Vars(r)
	targetName := vars["name"]
	conf, ok := gc.findEndpoint(targetName)
	if !ok {
		log.Printf("Could not find '%s'", targetName)
		return
	}
	gc.editConf(*conf, w, r)
}

func (gc apiGenericConfig[T]) editConf(conf T, w http.ResponseWriter, r *http.Request) {

	dat := tmplProviderGenericData[T]{
		System: system,
		Title:  fmt.Sprintf("Editing %s", gc.ConfTypeName),
		Conf:   conf,
	}
	err := templates.ExecuteTemplate(w, "Provider_Generic.html", dat)
	if err != nil {
		log.Printf("problem with template. %v", err)
	}
}

func (gc apiGenericConfig[T]) api_NewConf(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	var conf T
	system.Changes = true

	gc.Confs = append(gc.Confs, conf)

	gc.editConf(conf, w, r)
}

func (gc apiGenericConfig[T]) api_EditEndpoint(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	vars := mux.Vars(r)
	targetName := vars["name"]
	conf, ok := gc.findEndpoint(targetName)
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
	} else if index < 0 || index >= len(conf.Endpoints) {
		log.Printf("invalid endpoint %d.  must be 0.. %d", index, len(conf.Endpoints))
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
	newendpoint := conf.Endpoints[index]
	newendpoint.Historian = r.FormValue("Historian")
	newendpoint.Name = r.FormValue("Name")
	newendpoint.TagName = r.FormValue("TagName")
	newendpoint.Rate = rate
	newendpoint.TagType = gologix.CIPType(ciptype)
	conf.Endpoints[index] = newendpoint
	system.Changes = true

	gc.editConf(*conf, w, r)

}

func (gc apiGenericConfig[T]) api_NewEndpoint(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	vars := mux.Vars(r)
	targetName := vars["name"]
	conf, ok := gc.findEndpoint(targetName)
	log.Printf("Editing Endpoint %s", r.URL)
	if !ok {
		log.Printf("Could not find '%s'", targetName)
		return
	}

	system.Changes = true
	conf.Endpoints = append(conf.Endpoints, EndpointGeneric{})

	gc.editConf(*conf, w, r)

}

func (gc apiGenericConfig[T]) findEndpoint(name string) (*T, bool) {
	for i := range system.WorkingConfig.DataProviders.CIPClass3 {
		if system.WorkingConfig.DataProviders.CIPClass3[i].Name == name {
			return &gc.Confs[i], true
		}
	}

	log.Printf("Could not find '%s'", name)
	return new(T), false
}
