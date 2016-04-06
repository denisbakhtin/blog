package shared

import (
	"github.com/denisbakhtin/blog/models"
	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"log"
	"net/http"
)

var CSRF func(http.Handler) http.Handler

//InitCsrf initializes gorilla CSRF handler
func InitCsrf() {
	CSRF = csrf.Protect([]byte(GetConfig().CsrfSecret), csrf.Secure(GetConfig().Ssl), csrf.Path("/"), csrf.Domain(GetConfig().Domain))
}

//Default middleware chain for HandlerFuncs
func Default(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return CSRF(
		DataMiddleware(
			http.HandlerFunc(fn),
		),
	)
}

//Restricted middleware chain for HandlerFuncs
func Restricted(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return CSRF(
		RestrictedWithoutCSRF(fn),
	)
}

//RestrictedWithoutCSRF restriced middleware chain without CSRF protection
func RestrictedWithoutCSRF(fn func(http.ResponseWriter, *http.Request)) http.Handler {
	return DataMiddleware(
		RestrictedMiddleware(http.HandlerFunc(fn)),
	)
}

//DataMiddleware inits common request data (active user, et al). Must be preceded by SessionMiddleware
func DataMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		//set active user
		session := Session(r)
		if uID, ok := session.Values["user_id"]; ok {
			user, _ := models.GetUser(uID)
			if user.ID != 0 {
				context.Set(r, "user", user)
			}
		}
		//enable signup link
		if config.SignupEnabled {
			context.Set(r, "signup_enabled", true)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//RestrictedMiddleware verifies presence on 'user' in context (which is set by DataMiddleware, if user has signed in
func RestrictedMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if user := context.Get(r, "user"); user != nil {
			//access granted
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(403)
			tmpl.Lookup("errors/403").Execute(w, nil)
			log.Printf("ERROR: unauthorized access to %s\n", r.RequestURI)
		}
	}
	return http.HandlerFunc(fn)
}
