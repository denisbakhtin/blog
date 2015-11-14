package controllers

import (
	"fmt"
	"github.com/denisbakhtin/blog/helpers"
	"github.com/denisbakhtin/blog/models"
	"log"
	"net/http"
)

//Search handles POST /search route
func Search(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "POST" {

		query := r.PostFormValue("query")
		//full text search by name & description. Btw you can extend search to multi-table scenario with rankings, etc
		//fts index and SearchPosts assume language is english
		posts, _ := models.SearchPosts(query)
		data["Title"] = fmt.Sprintf("%s %q", T("search_results_for"), query)
		data["Posts"] = posts
		tmpl.Lookup("search/results").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}
