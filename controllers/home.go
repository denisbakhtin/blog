package controllers

import (
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
	"golang.org/x/net/context"
	"html/template"
)

//Home handles GET / route
func Home(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tmpl, data := ctx.Value("template").(*template.Template), helpers.DefaultData(ctx)
	if r.RequestURI != "/" {
		w.WriteHeader(404)
		tmpl.Lookup("errors/404").Execute(w, nil)
		return
	}
	data["Title"] = "Welcome to basic blog"
	data["Active"] = "home"
	tmpl.Lookup("home/show").Execute(w, data)
}
