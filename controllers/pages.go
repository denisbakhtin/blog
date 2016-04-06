package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/denisbakhtin/blog/models"
	"github.com/denisbakhtin/blog/shared"
)

//PageShow handles /pages/:id route
func PageShow(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	data := shared.DefaultData(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/pages/"):]
		page, err := models.GetPage(id)
		if err != nil || !page.Published {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Page"] = page
		data["Title"] = page.Name
		data["Active"] = fmt.Sprintf("pages/%s", id)
		tmpl.Lookup("pages/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}

//PageIndex handles GET /admin/pages route
func PageIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	data := shared.DefaultData(r)
	if r.Method == "GET" {

		list, err := models.GetPages()
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
			return
		}
		data["Title"] = "Pages"
		data["Active"] = "pages"
		data["List"] = list
		tmpl.Lookup("pages/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}

//PageCreate handles /admin/new_page route
func PageCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	session := shared.Session(r)
	data := shared.DefaultData(r)
	if r.Method == "GET" {

		data["Title"] = "New page"
		data["Active"] = "pages"
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("pages/form").Execute(w, data)

	} else if r.Method == "POST" {

		page := &models.Page{
			Name:      r.PostFormValue("name"),
			Content:   r.PostFormValue("content"),
			Published: shared.Atob(r.PostFormValue("published")),
		}

		if err := page.Insert(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, "/admin/new_page", 303)
			return
		}
		http.Redirect(w, r, "/admin/pages", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}

//PageUpdate handles /admin/edit_page/:id route
func PageUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	session := shared.Session(r)
	data := shared.DefaultData(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_page/"):]
		page, err := models.GetPage(id)
		if err != nil {
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
			return
		}

		data["Title"] = "Edit page"
		data["Active"] = "pages"
		data["Page"] = page
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("pages/form").Execute(w, data)

	} else if r.Method == "POST" {

		page := &models.Page{
			ID:        shared.Atoi64(r.PostFormValue("id")),
			Name:      r.PostFormValue("name"),
			Content:   r.PostFormValue("content"),
			Published: shared.Atob(r.PostFormValue("published")),
		}

		if err := page.Update(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/pages", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}

//PageDelete handles /admin/delete_page route
func PageDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)

	if r.Method == "POST" {

		page, err := models.GetPage(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, shared.ErrorData(err))
		}

		if err := page.Delete(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/pages", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}
