package controllers

import (
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
)

//Dashboard handles GET /admin route
func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	data["Title"] = T("dashboard")
	tmpl.Lookup("dashboard/show").Execute(w, data)
}
