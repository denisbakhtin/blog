package system

import (
	"github.com/denisbakhtin/blog/models"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
)

var (
	store *sessions.FilesystemStore
	//store *sessions.CookieStore
)

func createSession() {
	store = sessions.NewFilesystemStore("", []byte(config.SessionSecret))
	//store = sessions.NewCookieStore([]byte(config.SessionSecret))
	store.Options = &sessions.Options{Domain: config.Domain, Path: "/", Secure: config.Ssl, HttpOnly: true, MaxAge: 7 * 86400}
}

//SessionMiddleware creates gorilla session and stores it in context
func SessionMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer context.Clear(r)
		session, err := store.Get(r, "session")
		if err != nil {
			log.Printf("ERROR: can't get session: %s", err)
			http.Error(w, err.Error(), 500)
			return //abort chain
		}
		context.Set(r, "session", session)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//TemplateMiddleware stores parsed templates in context
func TemplateMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "template", tmpl)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

//DataMiddleware inits common request data (active user, et al). Must be preceded by SessionMiddleware
func DataMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session := context.Get(r, "session").(*sessions.Session)
		if uID, ok := session.Values["user_id"]; ok {
			user, _ := models.GetUser(uID)
			if user.ID != 0 {
				context.Set(r, "user", user)
			}
		}
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
			context.Get(r, "template").(*template.Template).Lookup("errors/403").Execute(w, nil)
			log.Printf("ERROR: unauthorized access to %s\n", r.RequestURI)
		}
	}
	return http.HandlerFunc(fn)
}
