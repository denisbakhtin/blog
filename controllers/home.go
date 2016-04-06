package controllers

import (
	"net/http"

	"github.com/denisbakhtin/blog/shared"
)

//Home handles GET / route
func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, data := shared.Template(r), shared.DefaultData(r)
	if r.RequestURI != "/" {
		w.WriteHeader(404)
		tmpl.Lookup("errors/404").Execute(w, nil)
		return
	}
	data["Title"] = "Welcome to blog boilerplate"
	data["Active"] = "home"
	tmpl.Lookup("home/show").Execute(w, data)
}
