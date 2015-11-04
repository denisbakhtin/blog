package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
	"github.com/denisbakhtin/blog/models"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"gopkg.in/guregu/null.v3"
)

//PostShow handles GET /posts/:id route
func PostShow(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)
	data := helpers.DefaultData(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/posts/"):]
		post, err := models.GetPost(id)
		if err != nil || !post.Published {
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Post"] = post
		data["Title"] = post.Name
		data["Active"] = fmt.Sprintf("posts/%s", id)
		tmpl.Lookup("posts/show").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PostIndex handles GET /admin/posts route
func PostIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)
	data := helpers.DefaultData(r)
	if r.Method == "GET" {

		list, err := models.GetPosts()
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Title"] = "List of posts"
		data["Active"] = "posts"
		data["List"] = list
		data[csrf.TemplateTag] = csrf.TemplateField(r)
		tmpl.Lookup("posts/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PostCreate handles /admin/new_post route
func PostCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)
	session := context.Get(r, "session").(*sessions.Session)
	data := helpers.DefaultData(r)
	if r.Method == "GET" {

		tags, err := models.GetTags()
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		data["Title"] = "Create post"
		data["Active"] = "posts"
		data["Tags"] = tags
		data["Flash"] = session.Flashes()
		data[csrf.TemplateTag] = csrf.TemplateField(r)
		session.Save(r, w)
		tmpl.Lookup("posts/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		post := &models.Post{
			Name:        r.PostFormValue("name"),
			Description: r.PostFormValue("description"),
			Published:   helpers.Atob(r.PostFormValue("published")),
			Tags:        r.Form["tags"], //PostFormValue returns only first value
		}

		if user := context.Get(r, "user"); user != nil {
			post.UserID = null.IntFrom(user.(*models.User).ID)
		}
		if err := post.Insert(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, "/admin/new_post", 303)
			return
		}
		http.Redirect(w, r, "/admin/posts", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PostUpdate handles /admin/edit_post/:id route
func PostUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)
	session := context.Get(r, "session").(*sessions.Session)
	data := helpers.DefaultData(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_post/"):]
		post, err := models.GetPost(id)
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}
		tags, err := models.GetTags()
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, nil)
			return
		}

		data["Title"] = "Edit post"
		data["Active"] = "posts"
		data["Post"] = post
		data["Tags"] = tags
		data["Flash"] = session.Flashes()
		data[csrf.TemplateTag] = csrf.TemplateField(r)
		session.Save(r, w)
		tmpl.Lookup("posts/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		post := &models.Post{
			ID:          helpers.Atoi64(r.PostFormValue("id")),
			Name:        r.PostFormValue("name"),
			Description: r.PostFormValue("description"),
			Published:   helpers.Atob(r.PostFormValue("published")),
			Tags:        r.Form["tags"], //PostFormValue returns only first value
		}

		if err := post.Update(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/posts", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//PostDelete handles /admin/delete_post route
func PostDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)

	if r.Method == "POST" {

		post, err := models.GetPost(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
		}

		if err := post.Delete(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/posts", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}
