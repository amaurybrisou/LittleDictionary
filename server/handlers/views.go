package handlers

import (
	// "log"
	"html/template"
	"net/http"
)


func IndexHandler(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("server/views/index.html")
	if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func AddHandler(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("server/views/add.html")
	if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func FindHandler(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("server/views/find.html")
	if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}
