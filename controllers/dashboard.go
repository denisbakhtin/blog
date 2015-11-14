package controllers

import (
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
	"github.com/gorilla/context"
	"github.com/nicksnyder/go-i18n/i18n"
	"html/template"
)

//Dashboard handles GET /admin route
func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)
	data := helpers.DefaultData(r)
	T := context.Get(r, "T").(i18n.TranslateFunc)
	data["Title"] = T("dashboard")
	tmpl.Lookup("dashboard/show").Execute(w, data)
}
