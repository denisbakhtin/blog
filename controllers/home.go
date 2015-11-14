package controllers

import (
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
)

//Home handles GET / route
func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, data := helpers.Template(r), helpers.DefaultData(r)
	T := helpers.T(r)
	if r.RequestURI != "/" {
		w.WriteHeader(404)
		tmpl.Lookup("errors/404").Execute(w, nil)
		return
	}
	data["Title"] = T("greeting")
	data["Active"] = "home"
	tmpl.Lookup("home/show").Execute(w, data)
}
