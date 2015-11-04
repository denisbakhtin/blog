package controllers

import (
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
	"github.com/gorilla/context"
	"html/template"
)

//Home handles GET / route
func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, data := context.Get(r, "template").(*template.Template), helpers.DefaultData(r)
	if r.RequestURI != "/" {
		w.WriteHeader(404)
		tmpl.Lookup("errors/404").Execute(w, nil)
		return
	}
	data["Title"] = "Welcome to basic blog"
	data["Active"] = "home"
	tmpl.Lookup("home/show").Execute(w, data)
}
