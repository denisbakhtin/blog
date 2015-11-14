package controllers

import (
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
	"github.com/gorilla/context"
	"github.com/nicksnyder/go-i18n/i18n"
	"html/template"
)

//Home handles GET / route
func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, data := context.Get(r, "template").(*template.Template), helpers.DefaultData(r)
	T := context.Get(r, "T").(i18n.TranslateFunc)
	if r.RequestURI != "/" {
		w.WriteHeader(404)
		tmpl.Lookup("errors/404").Execute(w, nil)
		return
	}
	data["Title"] = T("greeting")
	data["Active"] = "home"
	tmpl.Lookup("home/show").Execute(w, data)
}
