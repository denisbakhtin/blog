package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/denisbakhtin/blog/models"
	"github.com/denisbakhtin/blog/shared"
	"github.com/gorilla/context"
	"gopkg.in/guregu/null.v3"
)

//CommentIndex handles GET /admin/comments route
func CommentIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	data := shared.DefaultData(r)
	if r.Method == "GET" {

		list, err := models.GetComments()
		if err != nil {
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
			return
		}
		data["Title"] = "Comments"
		data["Active"] = "comments"
		data["List"] = list
		tmpl.Lookup("comments/index").Execute(w, data)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}

//CommentCreate handles /new_comment route
func CommentCreate(w http.ResponseWriter, r *http.Request) {
	session := shared.Session(r)
	tmpl := shared.Template(r)
	if r.Method == "POST" {

		if _, ok := session.Values["oauth_name"]; !ok {
			err := fmt.Errorf("You are not authorized to post comments.")
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(405)
			tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
			return
		}

		comment := &models.Comment{
			PostID:     shared.Atoi64(r.PostFormValue("post_id")),
			AuthorName: session.Values["oauth_name"].(string),
			Content:    r.PostFormValue("content"),
			Published:  false, //comments are published by admin via dashboard
		}

		if err := comment.Insert(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(400)
			tmpl.Lookup("errors/400").Execute(w, shared.ErrorData(err))
			return
		}
		session.AddFlash("Thank you! Your comment will be visible after approval.", "comments")
		session.Save(r, w)
		http.Redirect(w, r, fmt.Sprintf("/posts/%d#comments", comment.PostID), 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}

//CommentUpdate handles /admin/edit_comment/:id route
func CommentUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	session := shared.Session(r)
	data := shared.DefaultData(r)
	if r.Method == "GET" {

		id := r.URL.Path[len("/admin/edit_comment/"):]
		comment, err := models.GetComment(id)
		if err != nil {
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, shared.ErrorData(err))
			return
		}

		data["Title"] = "Edit comment"
		data["Active"] = "comments"
		data["Comment"] = comment
		data["Flash"] = session.Flashes("comments")
		session.Save(r, w)
		tmpl.Lookup("comments/form").Execute(w, data)

	} else if r.Method == "POST" {

		r.ParseForm()
		comment := &models.Comment{
			ID:        shared.Atoi64(r.PostFormValue("id")),
			Content:   r.PostFormValue("content"),
			Published: shared.Atob(r.PostFormValue("published")),
		}

		if err := comment.Update(); err != nil {
			session.AddFlash(err.Error(), "comments")
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/comments", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}

//CommentReply handles /admin/new_comment route
func CommentReply(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)
	session := shared.Session(r)
	data := shared.DefaultData(r)
	if r.Method == "GET" {

		user := context.Get(r, "user").(*models.User)
		parentID := shared.Atoi64(r.FormValue("parent_id"))
		parent, _ := models.GetComment(parentID)
		comment := &models.Comment{
			PostID:     parent.PostID,
			ParentID:   null.NewInt(parentID, parentID > 0),
			AuthorName: user.Name,
		}

		data["Title"] = "Reply"
		data["Active"] = "comments"
		data["Comment"] = comment
		data["Flash"] = session.Flashes("comments")
		session.Save(r, w)
		tmpl.Lookup("comments/form").Execute(w, data)

	} else if r.Method == "POST" {

		parentID := shared.Atoi64(r.PostFormValue("parent_id"))
		comment := &models.Comment{
			PostID:     shared.Atoi64(r.PostFormValue("post_id")),
			ParentID:   null.NewInt(parentID, parentID > 0),
			AuthorName: r.PostFormValue("author_name"),
			Content:    r.PostFormValue("content"),
			Published:  shared.Atob(r.PostFormValue("published")),
		}

		if err := comment.Insert(); err != nil {
			session.AddFlash(err.Error(), "comments")
			session.Save(r, w)
			http.Redirect(w, r, r.RequestURI, 303)
			return
		}
		http.Redirect(w, r, "/admin/comments", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}

//CommentDelete handles /admin/delete_comment route
func CommentDelete(w http.ResponseWriter, r *http.Request) {
	tmpl := shared.Template(r)

	if r.Method == "POST" {

		comment, err := models.GetComment(r.PostFormValue("id"))
		if err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(404)
			tmpl.Lookup("errors/404").Execute(w, shared.ErrorData(err))
		}

		if err := comment.Delete(); err != nil {
			log.Printf("ERROR: %s\n", err)
			w.WriteHeader(500)
			tmpl.Lookup("errors/500").Execute(w, shared.ErrorData(err))
			return
		}
		http.Redirect(w, r, "/admin/comments", 303)

	} else {
		err := fmt.Errorf("Method %q not allowed", r.Method)
		log.Printf("ERROR: %s\n", err)
		w.WriteHeader(405)
		tmpl.Lookup("errors/405").Execute(w, shared.ErrorData(err))
	}
}
