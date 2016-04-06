package controllers

import (
	"log"
	"net/http"

	"fmt"
	"github.com/denisbakhtin/blog/models"
	"github.com/denisbakhtin/blog/shared"
)

//UserIndex handles GET /admin/users route
func UserIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	data := shared.DefaultData(r)
	if r.Method == "GET" {

		list, err := models.GetUsers()
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/404").Execute(w, shared.ErrorData(err))
			return
		}
		data["Title"] = "Users"
		data["Active"] = "users"
		data["List"] = list
		tmpl.Lookup("users/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}

//UserCreate handles /admin/new_user route
func UserCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	session := shared.Session(r)
	data := shared.DefaultData(r)
	if r.Method == "GET" {

		data["Title"] = "New user"
		data["Active"] = "users"
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("users/form").Execute(w, data)

	} else if r.Method == "POST" {

		user := &models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}

		if err := user.HashPassword(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
			return
		}
		if err := user.Insert(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, "/admin/new_user", 303)
			return
		}
		http.Redirect(w, r, "/admin/users", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}

//UserUpdate handles /admin/edit_user/:id route
func UserUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	session := shared.Session(r)
	data := shared.DefaultData(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_user/"):]
		user, err := models.GetUser(id)
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, shared.ErrorData(err))
			return
		}

		data["Title"] = "Edit user"
		data["Active"] = "users"
		data["User"] = user
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("users/form").Execute(w, data)

	} else if r.Method == "POST" {

		user := &models.User{
			ID:       shared.Atoi64(r.PostFormValue("id")),
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}

		if err := user.HashPassword(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
			return
		}
		if err := user.Update(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/users", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}

//UserDelete handles /admin/delete_user route
func UserDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)

	if r.Method == "POST" {

		user, err := models.GetUser(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, shared.ErrorData(err))
		}

		if err := user.Delete(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/users", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}
