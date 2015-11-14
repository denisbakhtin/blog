package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/denisbakhtin/blog/helpers"
	"github.com/denisbakhtin/blog/models"
)

//ArchiveShow handles GET /archives/:year-:month route
func ArchiveShow(w http.ResponseWriter, r *http.Request) {
	tmpl := helpers.Template(r)
	data := helpers.DefaultData(r)
	T := helpers.T(r)
	if r.Method == "GET" {

		param := r.URL.Path[len("/archives/"):]
		ym := strings.Split(param, "-")
		if len(ym) != 2 {
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, nil)
			return
		}
		year, _ := strconv.Atoi(ym[0])
		month, _ := strconv.Atoi(ym[1])
		list, err := models.GetPostsByArchive(year, month)
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, nil)
			return
		}
		data["Title"] = fmt.Sprintf("%s %s %d", T("archives_for"), T(time.Month(month).String()), year)
		data["List"] = list
		data["Active"] = fmt.Sprintf("archives/%d-%02d", year, month)
		tmpl.Lookup("archives/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}
