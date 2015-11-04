package controllers

import (
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
	"golang.org/x/net/context"
	"html/template"
)

//Dashboard handles GET /admin route
func Dashboard(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tmpl, data := ctx.Value("template").(*template.Template), helpers.DefaultData(ctx)
	data["Title"] = "Blog dashboard"
	tmpl.Lookup("dashboard/show").Execute(w, data)
}
