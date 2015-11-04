package controllers

import (
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
	"github.com/gorilla/context"
	"html/template"
)

//Dashboard handles GET /admin route
func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)
	data := helpers.DefaultData(r)
	data["Title"] = "Blog dashboard"
	tmpl.Lookup("dashboard/show").Execute(w, data)
}
