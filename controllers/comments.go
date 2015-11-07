package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/denisbakhtin/blog/helpers"
	"github.com/denisbakhtin/blog/models"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"gopkg.in/guregu/null.v3"
)

//CommentIndex handles GET /admin/comments route
func CommentIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)
	data := helpers.DefaultData(r)
	if r.Method == "GET" {

		list, err := models.GetComments()
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		data["Title"] = "List of comments"
		data["Active"] = "comments"
		data["List"] = list
		tmpl.Lookup("comments/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//CommentCreate handles /new_comment route
func CommentCreate(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*sessions.Session)
	tmpl := context.Get(r, "template").(*template.Template)
	if r.Method == "POST" {

		parentID := helpers.Atoi64(r.PostFormValue("parent_id"))
		comment := &models.Comment{
			PostID:      helpers.Atoi64(r.PostFormValue("post_id")),
			ParentID:    null.NewInt(parentID, parentID > 0),
			AuthorName:  r.PostFormValue("author_name"), //TODO: get from cookies or session
			Description: r.PostFormValue("description"),
			Published:   false, //comments are published by admin via dashboard
		}

		if err := comment.Insert(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, helpers.ErrorData(err))
			return
		}
		session.AddFlash("Thank you! Your comment will be visible after approval.")
		session.Save(r, w)
		//TODO: show flash msg in comments block on that post page
		http.Redirect(w, r, fmt.Sprintf("/posts/%d#comments", comment.PostID), 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//CommentUpdate handles /admin/edit_comment/:id route
func CommentUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)
	session := context.Get(r, "session").(*sessions.Session)
	data := helpers.DefaultData(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_comment/"):]
		comment, err := models.GetComment(id)
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
			return
		}

		data["Title"] = "Edit comment"
		data["Active"] = "comments"
		data["Comment"] = comment
		data["Flash"] = session.Flashes()
		session.Save(r, w)
		tmpl.Lookup("comments/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		comment := &models.Comment{
			ID:          helpers.Atoi64(r.PostFormValue("id")),
			Description: r.PostFormValue("description"),
			Published:   helpers.Atob(r.PostFormValue("published")),
		}

		if err := comment.Update(); err != nil {
			session.AddFlash(err.Error())
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/comments", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}

//CommentDelete handles /admin/delete_comment route
func CommentDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := context.Get(r, "template").(*template.Template)

	if r.Method == "comment" {

		comment, err := models.GetComment(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, helpers.ErrorData(err))
		}

		if err := comment.Delete(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, helpers.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/comments", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, helpers.ErrorData(err))
	}
}
