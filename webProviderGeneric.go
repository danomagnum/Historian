package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

type DataProvider interface {
	Name() string
	String() string
	Update(url.Values) error
	Endpoints() []any
	NewEndpoint()
	UpdateEndpoint(url.Values) error
	RemoveEndpoint(int)
}

type ApiGenericConfig[T DataProvider] struct {
	ConfTypeName string
	Confs        []T
}

type tmplProviderGenericData struct {
	System System
	Title  string
	Conf   DataProvider
}

func (gc ApiGenericConfig[T]) Init(r *mux.Router) {
	r.HandleFunc("/{name}/Edit/", gc.api_EditConf)
	r.HandleFunc("/{name}/EditEndpoint/", gc.api_EditEndpoint)
	r.HandleFunc("/{name}/NewEndpoint/", gc.api_NewEndpoint)
	r.HandleFunc("/Add/", gc.api_NewConf)
}

func (gc ApiGenericConfig[T]) api_EditConf(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	vars := mux.Vars(r)
	targetName := vars["name"]
	conf, ok := gc.findConfByName(targetName)
	if !ok {
		log.Printf("Could not find '%s'", targetName)
		return
	}
	gc.editConf(*conf, w, r)
}

func (gc ApiGenericConfig[T]) editConf(conf DataProvider, w http.ResponseWriter, r *http.Request) {

	dat := tmplProviderGenericData{
		System: system,
		Title:  fmt.Sprintf("Editing %s", gc.ConfTypeName),
		Conf:   conf,
	}
	err := templates.ExecuteTemplate(w, "Provider_Generic.html", dat)
	if err != nil {
		log.Printf("problem with template. %v", err)
	}
}

func (gc ApiGenericConfig[T]) api_NewConf(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	var conf T
	system.Changes = true

	gc.Confs = append(gc.Confs, conf)

	gc.editConf(conf, w, r)
}

func (gc ApiGenericConfig[T]) api_EditEndpoint(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	vars := mux.Vars(r)
	targetName := vars["name"]
	conf, ok := gc.findConfByName(targetName)
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

	c := *conf
	endpoints := c.Endpoints()

	// see what we're editing.
	index_str := r.FormValue("Index")
	index, err := strconv.Atoi(index_str)
	if err != nil {
		log.Printf("invalid endpoint %s not an int: %v", index_str, err)
		return
	} else if index < 0 || index >= len(endpoints) {
		log.Printf("invalid endpoint %d.  must be 0.. %d", index, len(endpoints))
		return
	}

	err = c.UpdateEndpoint(r.Form)
	if err != nil {
		log.Printf("Could not update endpoint '%s': %v", targetName, err)
		return
	}

	gc.editConf(*conf, w, r)

}

func (gc ApiGenericConfig[T]) api_NewEndpoint(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob(templatedir + "*") // TODO: remove once page debug is done
	vars := mux.Vars(r)
	targetName := vars["name"]
	conf, ok := gc.findConfByName(targetName)
	log.Printf("Editing Endpoint %s", r.URL)
	if !ok {
		log.Printf("Could not find '%s'", targetName)
		return
	}

	system.Changes = true
	c := *conf
	c.NewEndpoint()

	gc.editConf(*conf, w, r)

}

func (gc ApiGenericConfig[T]) findConfByName(name string) (*T, bool) {
	for i := range gc.Confs {
		if gc.Confs[i].Name() == name {
			return &gc.Confs[i], true
		}
	}

	log.Printf("Could not find '%s'", name)
	return new(T), false
}
