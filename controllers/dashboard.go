package controllers

import (
	"net/http"

	"github.com/denisbakhtin/blog/shared"
)

//Dashboard handles GET /admin route
func Dashboard(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	data := shared.DefaultData(r)
	data["Title"] = "Dashboard"
	tmpl.Lookup("dashboard/show").Execute(w, data)
}
