package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type tmplProviderClass3Data struct {
	Title string
	Conf  ConfigCIPClass3
}

var webApiCIPClass3_Handler = http.NewServeMux()

func init() {
	webApiCIPClass3_Handler.HandleFunc("/Edit/", api_EditCipClass3Conf)
	webApiCIPClass3_Handler.HandleFunc("/Add/", api_NewCipClass3Conf)
}

func api_EditCipClass3Conf(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob("./templates/*") // TODO: remove once page debug is done
	targetName := strings.TrimPrefix(r.URL.Path, "/Edit/")
	var conf ConfigCIPClass3
	invalid := true
	for i := range workingConf.DataProviders.CIPClass3 {
		if workingConf.DataProviders.CIPClass3[i].Name == targetName {
			conf = workingConf.DataProviders.CIPClass3[i]
			invalid = false
			break
		}
	}

	if invalid {
		log.Printf("Could not find '%s'", targetName)
		return
	}
	dat := tmplProviderClass3Data{
		Title: fmt.Sprintf("Editing CIP Class 3 Endpoint %s @ %s,%s", conf.Name, conf.Address, conf.Path),
		Conf:  conf,
	}
	err := templates.ExecuteTemplate(w, "Provider_CIPClass3.html", dat)
	if err != nil {
		log.Printf("problem with template. %v", err)
	}
}

func api_NewCipClass3Conf(w http.ResponseWriter, r *http.Request) {
	templates, _ = template.ParseGlob("./templates/*") // TODO: remove once page debug is done
	dat := tmplProviderClass3Data{
		Title: "New CIP Class 3 Endpoint",
		Conf:  ConfigCIPClass3{},
	}
	err := templates.ExecuteTemplate(w, "Provider_CIPClass3.html", dat)
	if err != nil {
		log.Printf("problem with template. %v", err)
	}
}
